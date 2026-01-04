package grpc

import (
	"context"
	"log"
	"net"

	pb "market/market/proto"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MarketServer реализует MarketService
type MarketServer struct {
	pb.UnimplementedMarketServiceServer
}

// GetWeaponByName - получение оружия по имени
func (s *MarketServer) GetWeaponByName(ctx context.Context, req *pb.GetWeaponByNameRequest) (*pb.GetWeaponByNameResponse, error) {
	name := req.GetName()
	log.Printf("GRPC MarketServer: GetWeaponByName called: %s", name)

	// TODO: Подключите вашу реальную логику
	// weapon, err := weaponService.GetByName(ctx, name)
	// if err != nil {
	//     return nil, err
	// }

	// Заглушка
	weapon := &pb.Weapon{
		Id:            1,
		Name:          name,
		WeaponType:    "Sword",
		DamageMin:     10,
		DamageMax:     20,
		Weight:        5.5,
		Price:         1000,
		Rarity:        "Epic",
		SpecialEffect: "Fire damage",
		CreatedAt:     timestamppb.Now(),
	}

	return &pb.GetWeaponByNameResponse{
		Weapon: weapon,
	}, nil
}

// GetArmorByName - получение брони по имени
func (s *MarketServer) GetArmorByName(ctx context.Context, req *pb.GetArmorByNameRequest) (*pb.GetArmorByNameResponse, error) {
	name := req.GetName()
	log.Printf("GRPC MarketServer: GetArmorByName called: %s", name)

	// TODO: Подключите вашу реальную логику
	// armor, err := armorService.GetByName(ctx, name)
	// if err != nil {
	//     return nil, err
	// }

	// Заглушка
	armor := &pb.Armor{
		Id:            1,
		Name:          name,
		ArmorType:     "Plate",
		Defense:       50,
		Weight:        20.0,
		Price:         2000,
		Rarity:        "Rare",
		SpecialEffect: "Magic resistance",
		CreatedAt:     timestamppb.Now(),
	}

	return &pb.GetArmorByNameResponse{
		Armor: armor,
	}, nil
}

// StartServer запускает GRPC сервер
func StartServer(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	server := &MarketServer{}
	pb.RegisterMarketServiceServer(s, server)

	log.Printf("Market (Inventory) GRPC Server listening on :%s", port)
	return s.Serve(lis)
}
