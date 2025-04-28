package utils

import (
	"context"
	"demo-prismao-apicodegen/prisma/db"
	"fmt"
	"log"
)

// PrintUsers hiển thị thông tin về các user trong database
func PrintUsers(client *db.PrismaClient) {
	ctx := context.Background()

	// Lấy tất cả users từ database
	users, err := client.User.FindMany().Exec(ctx)
	if err != nil {
		log.Fatalf("Failed to query users: %v", err)
	}

	// In thông tin users
	if len(users) == 0 {
		fmt.Println("No users found in the database")
		return
	}

	fmt.Println("=== USER LIST ===")
	for i, user := range users {
		fmt.Printf("\nUser #%d:\n", i+1)
		fmt.Printf("ID: %d\n", user.ID)
		fmt.Printf("Name: %s\n", user.Name)
		fmt.Printf("Email: %s\n", user.Email)
		fmt.Printf("Created At: %s\n", user.CreatedAt)

		// Nếu muốn hiển thị cả roles của user (quan hệ User -> UserRoles -> Role)
		if roles, err := client.UserRole.FindMany(
			db.UserRole.UserID.Equals(user.ID),
		).With(
			db.UserRole.Role.Fetch(),
		).Exec(ctx); err == nil && len(roles) > 0 {
			fmt.Println("Roles:")
			for _, role := range roles {
				if role.Role() != nil {
					fmt.Printf("- %s (ID: %d)\n", role.Role().Name, role.RoleID)
				}
			}
		}
	}
	fmt.Println("\n=== END ===")
}
