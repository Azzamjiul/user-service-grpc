package main

import (
	"log"
	"net"
	"strconv"

	"user-service/domain/user"
	"user-service/domain/user/repository"
	"user-service/domain/user/service"
	"user-service/grpc/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Inisialisasi database
	dsn := "root:@tcp(127.0.0.1:3306)/user-service?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto Migrate
	if err := db.AutoMigrate(&user.User{}); err != nil {
		// Tangani kesalahan saat melakukan AutoMigrate
		panic("failed to auto migrate database: " + err.Error())
	}

	// Inisialisasi repository
	userRepository := repository.NewUserRepository(db)

	// Inisialisasi server Gin
	go func() {
		r := gin.Default()

		// Handler untuk membuat pengguna baru
		r.POST("/users", func(c *gin.Context) {
			var newUser user.User
			if err := c.ShouldBindJSON(&newUser); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			if err := userRepository.Create(&newUser); err != nil {
				c.JSON(500, gin.H{"error": "failed to create user"})
				return
			}
			c.JSON(201, newUser)
		})

		// Handler untuk mendapatkan pengguna berdasarkan ID
		r.GET("/users/:id", func(c *gin.Context) {
			userID := c.Param("id")
			userIDInt, err := strconv.Atoi(userID)
			if err != nil {
				c.JSON(400, gin.H{"error": "id must be integer"})
				return
			}
			user, err := userRepository.FindByID(uint(userIDInt))
			if err != nil {
				c.JSON(404, gin.H{"error": "user not found"})
				return
			}
			c.JSON(200, user)
		})

		// Jalankan server Gin
		if err := r.Run(":8080"); err != nil {
			panic("failed to start Gin server: " + err.Error())
		}
	}()

	// Inisialisasi server gRPC
	go func() {
		// Inisialisasi server gRPC
		grpcServer := grpc.NewServer()

		// Daftarkan implementasi layanan gRPC
		// proto.RegisterUserServiceServer(grpcServer, &service.UserServiceServer{})
		proto.RegisterUserServiceServer(grpcServer, &service.UserServiceServer{UserRepository: *userRepository})

		// Mulai server gRPC
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Tunggu tanpa melakukan apa pun
	select {}
}
