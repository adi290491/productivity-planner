package util

import "testing"

func BenchmarkHashPassword(b *testing.B) {
	password := "testpassword123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HashPassword(password)
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	password := "testpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		b.Fatalf("Failed to hash password: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = VerifyPassword(password, hash)
	}
}

func BenchmarkHashPassword_Parallel(b *testing.B) {
	password := "testpassword123"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = HashPassword(password)
		}
	})
}

func BenchmarkVerifyPassword_Parallel(b *testing.B) {
	password := "testpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		b.Fatalf("Failed to hash password: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = VerifyPassword(password, hash)
		}
	})
}
