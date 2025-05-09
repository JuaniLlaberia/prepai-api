package configs

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var DB *mongo.Client
var Ctx context.Context
var Cancel context.CancelFunc

func ConnectDB() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(ProcessEnv("MONGOURI")).SetServerAPIOptions(serverAPI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	DB, err = mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	if err := DB.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Database connected successfully")
}

func DisconnectDB() {
	if DB != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := DB.Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect from database: %v", err)
		} else {
			fmt.Println("Database disconnected successfully")
		}
	}
}

func GetCollection(collectionName string) *mongo.Collection {
	return DB.Database("PrepAi").Collection(collectionName)
}

func InitDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	collections := []struct {
		name  string
		setup func(context.Context) error
	}{
		{"users", SetupUserCollection},
		{"interviews", SetupInterviewCollection},
		{"exams", SetupExamCollection},
		{"questions", SetupQuestionCollection},
		{"resumes", SetupResumeCollection},
		{"interviewAttempts", SetupInterviewAttemptCollection},
		{"examAttempts", SetupExamAttemptCollection},
	}

	for _, col := range collections {
		wg.Add(1)
		go func(name string, setup func(context.Context) error) {
			defer wg.Done()
			if err := col.setup(ctx); err != nil {
				errChan <- fmt.Errorf("failed to set up %v collection: %v", col.name, err)
			}
		}(col.name, col.setup)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		log.Fatalf("%v", err)
	}
}

func SetupUserCollection(ctx context.Context) error {
	collection := GetCollection("users")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"email", 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("failed to create email index: %v", err)
	}

	jsonSchema := bson.M{
		"bsonType": "object",
		"required": []string{"email", "password"},
		"properties": bson.M{
			"full_name": bson.M{
				"bsonType":    "string",
				"description": "Full name of the user",
			},
			"email": bson.M{
				"bsonType":    "string",
				"pattern":     "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
				"description": "Email must be a valid email address",
			},
			"password": bson.M{
				"bsonType":    "string",
				"minLength":   8,
				"description": "Password must be at least 8 characters",
			},
			"image_url": bson.M{
				"bsonType":    "string",
				"description": "URL to user's profile image",
			},
		},
	}

	validator := bson.M{
		"$jsonSchema": jsonSchema,
	}

	command := bson.D{
		{"collMod", "users"},
		{"validator", validator},
		{"validationLevel", "moderate"},
	}

	err = DB.Database("PrepAi").RunCommand(ctx, command).Err()
	if err != nil {
		if strings.Contains(err.Error(), "namespace") {
			createOpts := options.CreateCollection().SetValidator(validator)
			err = DB.Database("PrepAi").CreateCollection(ctx, "users", createOpts)
			if err != nil {
				return fmt.Errorf("failed to create users collection: %v", err)
			}
		} else {
			return fmt.Errorf("failed to set up validator: %v", err)
		}
	}

	return nil
}

func SetupInterviewCollection(ctx context.Context) error {
	collection := GetCollection("interviews")

	_, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{"user_id", 1}},
		},
		{
			Keys: bson.D{{"activity_id", 1}},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create user_id/activity_id index: %v", err)
	}

	jsonSchema := bson.M{
		"bsonType": "object",
		"required": []string{"job_role", "job_level", "topics"},
		"properties": bson.M{
			"title": bson.M{
				"bsonType":    "string",
				"description": "Interview descriptive title",
			},
			"job_role": bson.M{
				"bsonType":    "string",
				"description": "Job name needed for the interview",
			},
			"job_level": bson.M{
				"bsonType":    "string",
				"enum":        []string{"intership", "junior", "ssr", "senior", "lead"},
				"description": "Interview difficulty based on seniority",
			},
			"topics": bson.M{
				"bsonType": "array",
				"items": bson.M{
					"description": "List of topics",
					"bsonType":    "string",
				},
				"minItems":    1,
				"uniqueItems": true,
			},
			"taken": bson.M{
				"bsonType":    "bool",
				"description": "Describes if the interview was taken by the user or not",
			},
			"pinned": bson.M{
				"bsonType":    "bool",
				"description": "Describes if the user pinned to top the interview",
			},
			"passed": bson.M{
				"bsonType":    "bool",
				"description": "Describes if the user passed or not the interview",
			},
			"questions": bson.M{
				"type": "array",
				"items": bson.M{
					"description": "Interview questions",
					"bsonType":    "object",
					"properties": bson.M{
						"question": bson.M{
							"bsonType":    "string",
							"description": "Interview question",
						},
						"hint": bson.M{
							"bsonType":    "string",
							"description": "Helper text for user",
						},
						"type": bson.M{
							"bsonType":    "string",
							"description": "Question type",
						},
					},
				},
				"minItems":    1,
				"uniqueItems": true,
			},
			"user_id": bson.M{
				"bsonType":    "objectId",
				"description": "Reference to user who created the interview",
			},
			"activity_id": bson.M{
				"bsonType":    "objectId",
				"description": "Reference to activity for which this interview belongs to",
			},
		},
	}

	validator := bson.M{
		"$jsonSchema": jsonSchema,
	}

	command := bson.D{
		{"collMod", "interviews"},
		{"validator", validator},
		{"validationLevel", "moderate"},
	}

	err = DB.Database("PrepAi").RunCommand(ctx, command).Err()
	if err != nil {
		if strings.Contains(err.Error(), "namespace") {
			createOpts := options.CreateCollection().SetValidator(validator)
			err = DB.Database("PrepAi").CreateCollection(ctx, "interviews", createOpts)
			if err != nil {
				return fmt.Errorf("failed to create interviews collection: %v", err)
			}
		} else {
			return fmt.Errorf("failed to set up validator: %v", err)
		}
	}

	return nil
}

func SetupExamCollection(ctx context.Context) error {
	collection := GetCollection("exams")

	_, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{"user_id", 1}},
		},
		{
			Keys: bson.D{{"activity_id", 1}},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create user_id/activity_id index: %v", err)
	}

	jsonSchema := bson.M{
		"bsonType": "object",
		"required": []string{"subject", "difficulty", "type"},
		"properties": bson.M{
			"title": bson.M{
				"bsonType":    "string",
				"description": "Exam descriptive title",
			},
			"subject": bson.M{
				"bsonType":    "string",
				"description": "Exam topic",
			},
			"difficulty": bson.M{
				"bsonType":    "string",
				"enum":        []string{"easy", "medium", "hard"},
				"description": "How easy/hard is the exam",
			},
			"exam_type": bson.M{
				"bsonType":    "string",
				"enum":        []string{"true-false", "multiple-choice"},
				"description": "Based on the type the answers will change",
			},
			"taken": bson.M{
				"bsonType":    "bool",
				"description": "Describes if the exam was taken by the user or not",
			},
			"pinned": bson.M{
				"bsonType":    "bool",
				"description": "Describes if the user pinned to top the exam",
			},
			"passed": bson.M{
				"bsonType":    "bool",
				"description": "Describes if the user passed or not the exam",
			},
			"questions": bson.M{
				"type": "array",
				"items": bson.M{
					"description": "Exam questions",
					"bsonType":    "object",
					"properties": bson.M{
						"question": bson.M{
							"bsonType":    "string",
							"description": "Exam question",
						},
						"options": bson.M{
							"bsonType": "array",
							"items": bson.M{
								"description": "Options to answer question",
								"bsonType":    "string",
							},
							"minItems":    1,
							"uniqueItems": true,
						},
						"correct": bson.M{
							"bsonType":    "number",
							"description": "Question correct answer (index)",
						},
					},
				},
				"minItems":    1,
				"uniqueItems": true,
			},
			"user_id": bson.M{
				"bsonType":    "objectId",
				"description": "Reference to user who created the exam",
			},
			"activity_id": bson.M{
				"bsonType":    "objectId",
				"description": "Reference to activity for which this exam belongs to",
			},
		},
	}

	validator := bson.M{
		"$jsonSchema": jsonSchema,
	}

	command := bson.D{
		{"collMod", "exams"},
		{"validator", validator},
		{"validationLevel", "moderate"},
	}

	err = DB.Database("PrepAi").RunCommand(ctx, command).Err()
	if err != nil {
		if strings.Contains(err.Error(), "namespace") {
			createOpts := options.CreateCollection().SetValidator(validator)
			err = DB.Database("PrepAi").CreateCollection(ctx, "exams", createOpts)
			if err != nil {
				return fmt.Errorf("failed to create exams collection: %v", err)
			}
		} else {
			return fmt.Errorf("failed to set up validator: %v", err)
		}
	}

	return nil
}

func SetupQuestionCollection(ctx context.Context) error {
	collection := GetCollection("questions")

	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{"user_id", 1}},
	})
	if err != nil {
		return fmt.Errorf("failed to create user_id index: %v", err)
	}

	jsonSchema := bson.M{
		"bsonType": "object",
		"required": []string{"question", "user_id"},
		"properties": bson.M{
			"question": bson.M{
				"bsonType":    "string",
				"description": "Question to analyse",
			},
			"type": bson.M{
				"bsonType":    "string",
				"description": "Question type",
			},
			"difficulty": bson.M{
				"bsonType":    "string",
				"enum":        []string{"easy", "medium", "hard"},
				"description": "How easy/hard is the exam",
			},
			"explanation": bson.M{
				"bsonType":    "string",
				"description": "Question explanation",
			},
			"expected_length": bson.M{
				"bsonType":    "string",
				"description": "Describes how much time should someone take to answer this question",
			},
			"pinned": bson.M{
				"bsonType":    "bool",
				"description": "Describes if the user pinned to top the exam",
			},
			"ideal_answer": bson.M{
				"description": "Object containing data on how to answer the question",
				"bsonType":    "object",
				"properties": bson.M{
					"structure": bson.M{
						"bsonType": "string",
					},
					"key_points": bson.M{
						"bsonType": "array",
						"items": bson.M{
							"bsonType": "string",
						},
						"minItems":    1,
						"uniqueItems": true,
					},
					"example": bson.M{
						"bsonType":    "string",
						"description": "An answered generated by AI on how to answer this question",
					},
				},
				"minItems":    1,
				"uniqueItems": true,
			},
			"user_id": bson.M{
				"bsonType":    "objectId",
				"description": "Reference to user who created the question analysis",
			},
		},
	}

	validator := bson.M{
		"$jsonSchema": jsonSchema,
	}

	command := bson.D{
		{"collMod", "questions"},
		{"validator", validator},
		{"validationLevel", "moderate"},
	}

	err = DB.Database("PrepAi").RunCommand(ctx, command).Err()
	if err != nil {
		if strings.Contains(err.Error(), "namespace") {
			createOpts := options.CreateCollection().SetValidator(validator)
			err = DB.Database("PrepAi").CreateCollection(ctx, "questions", createOpts)
			if err != nil {
				return fmt.Errorf("failed to create questions collection: %v", err)
			}
		} else {
			return fmt.Errorf("failed to set up validator: %v", err)
		}
	}

	return nil
}

func SetupResumeCollection(ctx context.Context) error {
	collection := GetCollection("resumes")

	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{"user_id", 1}},
	})
	if err != nil {
		return fmt.Errorf("failed to create user_id index: %v", err)
	}

	jsonSchema := bson.M{
		"bsonType": "object",
		"required": []string{"user_id"},
		"properties": bson.M{
			"file_url": bson.M{
				"bsonType":    "string",
				"description": "The url where the resume is stored in the cloud",
			},
			"title": bson.M{
				"bsonType":    "string",
				"description": "Descriptive title of this resume analysis",
			},
			"overall_score": bson.M{
				"bsonType":    "number",
				"description": "How well did the resume performed",
			},
			"analysis_summery": bson.M{
				"bsonType":    "string",
				"description": "Written analysis",
			},
			"improvement_suggestions": bson.M{
				"bsonType":    "string",
				"description": "List of suggestion to improve the resume",
			},
			"metrics": bson.M{
				"description": "Object containing data on how did the resume scored on certain categories",
				"bsonType":    "object",
				"properties": bson.M{
					"ats_match_score": bson.M{
						"bsonType": "number",
					},
					"clarity_score": bson.M{
						"bsonType": "number",
					},
					"grammar_issues": bson.M{
						"bsonType": "number",
					},
					"soft_vs_hard_skill_balance": bson.M{
						"bsonType": "string",
					},
					"resume_length_feedback": bson.M{
						"bsonType": "string",
					},
					"filler_word_usage": bson.M{
						"bsonType": "string",
					},
				},
				"minItems":    1,
				"uniqueItems": true,
			},
			"user_id": bson.M{
				"bsonType":    "objectId",
				"description": "Reference to user who created the resume analysis",
			},
		},
	}

	validator := bson.M{
		"$jsonSchema": jsonSchema,
	}

	command := bson.D{
		{"collMod", "resumes"},
		{"validator", validator},
		{"validationLevel", "moderate"},
	}

	err = DB.Database("PrepAi").RunCommand(ctx, command).Err()
	if err != nil {
		if strings.Contains(err.Error(), "namespace") {
			createOpts := options.CreateCollection().SetValidator(validator)
			err = DB.Database("PrepAi").CreateCollection(ctx, "resumes", createOpts)
			if err != nil {
				return fmt.Errorf("failed to create resumes collection: %v", err)
			}
		} else {
			return fmt.Errorf("failed to set up validator: %v", err)
		}
	}

	return nil
}

func SetupInterviewAttemptCollection(ctx context.Context) error {
	collection := GetCollection("interviewAttempts")

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "interview_id", Value: 1},
		},
		Options: options.Index().SetName("compoundIndex"),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("failed to create compoundIndex index: %v", err)
	}

	jsonSchema := bson.M{
		"bsonType": "object",
		"required": []string{"user_id", "interview_id", "answers"},
		"properties": bson.M{
			"answers": bson.M{
				"type": "array",
				"items": bson.M{
					"description": "Interview questions feedback and answers",
					"bsonType":    "object",
					"properties": bson.M{
						"question": bson.M{
							"bsonType":    "string",
							"description": "Interview question",
						},
						"user_response": bson.M{
							"bsonType":    "string",
							"description": "The response that the user gave for this question",
						},
						"feedback": bson.M{
							"bsonType":    "string",
							"description": "Feedback provided by ai based on user response",
						},
						"score": bson.M{
							"bsonType":    "number",
							"description": "How well did the user answered the question",
						},
						"suggestion": bson.M{
							"bsonType":    "string",
							"description": "How to improve the response",
						},
					},
				},
				"minItems":    1,
				"uniqueItems": true,
			},
			"analysis": bson.M{
				"bsonType":    "string",
				"description": "Speech analysis on the whole interview attempt",
			},
			"strengths": bson.M{
				"bsonType": "array",
				"items": bson.M{
					"description": "What are the user strengths in this interview",
					"bsonType":    "string",
				},
				"minItems":    1,
				"uniqueItems": true,
			},
			"areas_to_improve": bson.M{
				"bsonType": "array",
				"items": bson.M{
					"description": "What the user needs to improve",
					"bsonType":    "string",
				},
				"minItems":    1,
				"uniqueItems": true,
			},

			"passed": bson.M{
				"bsonType":    "bool",
				"description": "Describes if the user passed or not the interview",
			},
			"user_id": bson.M{
				"bsonType":    "objectId",
				"description": "Reference to user who created the resume analysis",
			},
			"interview_id": bson.M{
				"bsonType":    "objectId",
				"description": "Reference to the interview that belongs to",
			},
		},
	}

	validator := bson.M{
		"$jsonSchema": jsonSchema,
	}

	command := bson.D{
		{"collMod", "interviewAttempts"},
		{"validator", validator},
		{"validationLevel", "moderate"},
	}

	err = DB.Database("PrepAi").RunCommand(ctx, command).Err()
	if err != nil {
		if strings.Contains(err.Error(), "namespace") {
			createOpts := options.CreateCollection().SetValidator(validator)
			err = DB.Database("PrepAi").CreateCollection(ctx, "interviewAttempts", createOpts)
			if err != nil {
				return fmt.Errorf("failed to create interviewAttempts collection: %v", err)
			}
		} else {
			return fmt.Errorf("failed to set up validator: %v", err)
		}
	}

	return nil
}

func SetupExamAttemptCollection(ctx context.Context) error {
	collection := GetCollection("examAttempts")

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "exam_id", Value: 1},
		},
		Options: options.Index().SetName("compoundIndex"),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("failed to create compoundIndex index: %v", err)
	}

	jsonSchema := bson.M{
		"bsonType": "object",
		"required": []string{"user_id", "exam_id", "answers"},
		"properties": bson.M{
			"answers": bson.M{
				"type": "array",
				"items": bson.M{
					"description": "Exam questions feedback and answers",
					"bsonType":    "object",
					"properties": bson.M{
						"question": bson.M{
							"bsonType":    "string",
							"description": "Exam question",
						},
						"answer": bson.M{
							"bsonType":    "number",
							"description": "User answer (index)",
						},
						"correct": bson.M{
							"bsonType":    "number",
							"description": "Correct answer (index)",
						},
						"Explanation": bson.M{
							"bsonType":    "string",
							"description": "Brief explanation on why the correct answer is correct",
						},
					},
				},
				"minItems":    1,
				"uniqueItems": true,
			},
			"time": bson.M{
				"bsonType":    "number",
				"description": "The time that the user took to answer all questions in seconds",
			},
			"score": bson.M{
				"bsonType":    "number",
				"description": "Correct answers amount",
			},
			"passed": bson.M{
				"bsonType":    "bool",
				"description": "Describes if the user passed or not the interview",
			},
			"user_id": bson.M{
				"bsonType":    "objectId",
				"description": "Reference to user who created the resume analysis",
			},
			"exam_id": bson.M{
				"bsonType":    "objectId",
				"description": "Reference to the exam that belongs to",
			},
		},
	}

	validator := bson.M{
		"$jsonSchema": jsonSchema,
	}

	command := bson.D{
		{"collMod", "examAttempts"},
		{"validator", validator},
		{"validationLevel", "moderate"},
	}

	err = DB.Database("PrepAi").RunCommand(ctx, command).Err()
	if err != nil {
		if strings.Contains(err.Error(), "namespace") {
			createOpts := options.CreateCollection().SetValidator(validator)
			err = DB.Database("PrepAi").CreateCollection(ctx, "examAttempts", createOpts)
			if err != nil {
				return fmt.Errorf("failed to create examAttempts collection: %v", err)
			}
		} else {
			return fmt.Errorf("failed to set up validator: %v", err)
		}
	}

	return nil
}
