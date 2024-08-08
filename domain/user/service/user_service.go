package service

import (
	"context"
	"errors"
	"log"
	"user-service/domain/user/repository"
	"user-service/grpc/proto"
)

// Implementasi layanan gRPC
type UserServiceServer struct {
	UserRepository repository.UserRepository
}

// Implementasi metode GetUserById
func (s *UserServiceServer) GetUserById(ctx context.Context, req *proto.GetUserByIdRequest) (*proto.User, error) {
	// Ambil data pengguna dari database berdasarkan ID yang diberikan
	log.Println("Ambil data pengguna dari database berdasarkan ID yang diberikan")
	userEntity, err := s.UserRepository.FindByID(uint(req.UserId))
	if err != nil {
		// Tangani kesalahan saat mengambil data pengguna dari database
		// Anda dapat menyesuaikan pesan kesalahan berdasarkan jenis kesalahan yang diterima
		log.Println("failed to get user: " + err.Error())
		return nil, errors.New("failed to get user: " + err.Error())
	}

	log.Println("Ubah entitas pengguna ke tipe protobuf", userEntity)
	// Ubah entitas pengguna ke tipe protobuf
	userProto := &proto.User{
		Id:   uint64(userEntity.ID),
		Name: userEntity.Name,
		// Tambahkan bidang lain sesuai kebutuhan
	}

	return userProto, nil
}
