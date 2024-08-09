package main

import (
	"net"
	"os"
	"strconv"
	"time"

	"user-service/domain/user"
	"user-service/domain/user/repository"
	"user-service/domain/user/service"
	"user-service/grpc/proto"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var log = logrus.New()

func main() {
	// Set log format to JSON, which is useful for logging in Kubernetes
	log.SetFormatter(&logrus.JSONFormatter{})

	// Optionally set the log level
	log.SetLevel(logrus.InfoLevel)

	// Load configuration
	dsn := os.Getenv("DB_DSN")
	// dsn = "root:@tcp(localhost:3306)/user-service?charset=utf8mb4&parseTime=True&loc=Local"
	log.Println("dsn: " + dsn)

	// Initialize database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto Migrate
	if err := db.AutoMigrate(&user.User{}); err != nil {
		log.Fatalf("failed to auto migrate database: %v", err)
	}

	// Initialize repository
	userRepository := repository.NewUserRepository(db)

	// Initialize HTTP server
	go func() {
		r := gin.Default()

		r.POST("/users", func(c *gin.Context) {
			var newUser user.User
			if err := c.ShouldBindJSON(&newUser); err != nil {
				log.WithFields(logrus.Fields{
					"error": err.Error(),
				}).Error("Failed to bind JSON")
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			log.WithFields(logrus.Fields{
				"user": newUser,
			}).Info("New user created")

			if err := userRepository.Create(&newUser); err != nil {
				log.WithFields(logrus.Fields{
					"error": err.Error(),
				}).Error("Failed to create user")
				c.JSON(500, gin.H{"error": "failed to create user"})
				return
			}

			log.WithFields(logrus.Fields{
				"user": newUser,
			}).Info("User successfully created")
			c.JSON(201, newUser)
		})

		r.GET("/users/:id", func(c *gin.Context) {
			userID := c.Param("id")
			log.WithFields(logrus.Fields{
				"userID": userID,
			}).Info("Fetching user by ID")

			userIDInt, err := strconv.Atoi(userID)
			if err != nil {
				log.WithFields(logrus.Fields{
					"error":  "id must be integer",
					"userID": userID,
				}).Error("Failed to convert ID to integer")
				c.JSON(400, gin.H{"error": "id must be integer"})
				return
			}

			user, err := userRepository.FindByID(uint(userIDInt))
			if err != nil {
				log.WithFields(logrus.Fields{
					"error":  err.Error(),
					"userID": userID,
				}).Error("User not found")
				c.JSON(404, gin.H{"error": "user not found"})
				return
			}
			log.WithFields(logrus.Fields{
				"user": user,
			}).Info("User found")

			c.JSON(200, user)
		})

		if err := r.Run(":8080"); err != nil {
			log.Fatalf("failed to start Gin server: %v", err)
		}
	}()

	// Initialize gRPC server
	go func() {
		grpcServer := grpc.NewServer()

		proto.RegisterUserServiceServer(grpcServer, &service.UserServiceServer{UserRepository: *userRepository})

		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Graceful shutdown handling
	stop := make(chan struct{})
	go func() {
		<-stop
		// Add cleanup and graceful shutdown code here
		time.Sleep(2 * time.Second) // Simulate shutdown delay
		log.Println("Shutting down")
		os.Exit(0)
	}()

	// Block main thread
	select {}
}
