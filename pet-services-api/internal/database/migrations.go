package database

import (
	"pet-services-api/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Migration20260110000000(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Species{},
		&models.Breed{},
		&models.Category{},
		&models.Tag{},
		&models.Photo{},

		&models.Provider{},
		&models.Pet{},
		&models.Service{},
		&models.Review{},
		&models.Request{},
	)
}

func Migration20260321000000(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.AdoptionGuardianProfile{},
	)
}

func Migration20260215000000(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.RefreshToken{},
	)
}

func Migration20260204000000(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.PasswordResetToken{},
	)
}

func Migration20260213000000(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
	)
}

func Migration20260218000000(db *gorm.DB) error {
	const seed = `
INSERT INTO public.species
(id, "name", active, created_at, updated_at, deactivated_at, deleted_at)
VALUES
('01KG7BG1XN0BQ4KHKPHY5V5ZEW', 'Cachorro', true, '2026-01-30 08:44:27.573', '2026-01-30 08:44:27.575', NULL, NULL),
('01KG7BG1XR2FA63J6N46CPVCKS', 'Gato', true, '2026-01-30 08:44:27.576', '2026-01-30 08:44:27.576', NULL, NULL),
('01KG7BG1XR2FA63J6N4ACP7V5T', 'Pássaro', true, '2026-01-30 08:44:27.576', '2026-01-30 08:44:27.577', NULL, NULL),
('01KG7BG1XSGSDDR8NT15VN30PP', 'Peixe', true, '2026-01-30 08:44:27.577', '2026-01-30 08:44:27.577', NULL, NULL),
('01KG7BG1XSGSDDR8NT19E83CRF', 'Coelho', true, '2026-01-30 08:44:27.577', '2026-01-30 08:44:27.577', NULL, NULL),
('01KG7BG1XTY554EHAK9T2TQZYW', 'Hamster', true, '2026-01-30 08:44:27.578', '2026-01-30 08:44:27.578', NULL, NULL),
('01KG7BG1XTY554EHAK9X2WV61H', 'Tartaruga', true, '2026-01-30 08:44:27.578', '2026-01-30 08:44:27.578', NULL, NULL),
('01KG7BG1XTY554EHAK9ZEVPA57', 'Porquinho-da-índia', true, '2026-01-30 08:44:27.578', '2026-01-30 08:44:27.578', NULL, NULL),
('01KG7BG1XVSQR05Y3TPTNF09MX', 'Furão', true, '2026-01-30 08:44:27.579', '2026-01-30 08:44:27.579', NULL, NULL),
('01KG7BG1XVSQR05Y3TPTWZG03N', 'Cavalo', true, '2026-01-30 08:44:27.579', '2026-01-30 08:44:27.579', NULL, NULL),
('01KG7BG1XVSQR05Y3TPWMHJGKR', 'Pônei', true, '2026-01-30 08:44:27.579', '2026-01-30 08:44:27.580', NULL, NULL),
('01KG7BG1XW781NZ55BRAY9PKB5', 'Rato', true, '2026-01-30 08:44:27.580', '2026-01-30 08:44:27.580', NULL, NULL),
('01KG7BG1XW781NZ55BRE2TRBTW', 'Répteis', true, '2026-01-30 08:44:27.580', '2026-01-30 08:44:27.580', NULL, NULL),
('01KG7BG1XW781NZ55BRHE4YKQN', 'Anfíbios', true, '2026-01-30 08:44:27.580', '2026-01-30 08:44:27.581', NULL, NULL),
('01KG7BG1XXQ00TAASZYFX4V85R', 'Aranha', true, '2026-01-30 08:44:27.581', '2026-01-30 08:44:27.581', NULL, NULL),
('01KG7BG1XXQ00TAASZYH175BP8', 'Lagarto', true, '2026-01-30 08:44:27.581', '2026-01-30 08:44:27.581', NULL, NULL),
('01KG7BG1XYSCPQT8G8VENXRZ03', 'Papagaio', true, '2026-01-30 08:44:27.582', '2026-01-30 08:44:27.582', NULL, NULL),
('01KG7BG1XYSCPQT8G8VJ2VRJVW', 'Calopsita', true, '2026-01-30 08:44:27.582', '2026-01-30 08:44:27.582', NULL, NULL),
('01KG7BG1XYSCPQT8G8VKV5978X', 'Periquito', true, '2026-01-30 08:44:27.582', '2026-01-30 08:44:27.582', NULL, NULL),
('01KG7BG1XZ0GQ4CAGN2MJREQMB', 'Canário', true, '2026-01-30 08:44:27.583', '2026-01-30 08:44:27.583', NULL, NULL),
('01KG7BG1XZ0GQ4CAGN2RASVFH2', 'Cobra', true, '2026-01-30 08:44:27.583', '2026-01-30 08:44:27.583', NULL, NULL),
('01KG7BG1XZ0GQ4CAGN2SRQD5MD', 'Iguana', true, '2026-01-30 08:44:27.583', '2026-01-30 08:44:27.583', NULL, NULL),
('01KG7BG1Y0KVNE77CGEQR49KJ7', 'Chinchila', true, '2026-01-30 08:44:27.584', '2026-01-30 08:44:27.584', NULL, NULL),
('01KG7BG1Y0KVNE77CGEVJQEEPQ', 'Ouriço', true, '2026-01-30 08:44:27.584', '2026-01-30 08:44:27.584', NULL, NULL),
('01KG7BG1Y0KVNE77CGEXZEGZV7', 'Pato', true, '2026-01-30 08:44:27.584', '2026-01-30 08:44:27.584', NULL, NULL),
('01KG7BG1Y1CV582QT0675Z183A', 'Ganso', true, '2026-01-30 08:44:27.585', '2026-01-30 08:44:27.585', NULL, NULL),
('01KG7BG1Y1CV582QT068EYCPGJ', 'Galinha', true, '2026-01-30 08:44:27.585', '2026-01-30 08:44:27.585', NULL, NULL),
('01KG7BG1Y1CV582QT06AANRDQG', 'Pavão', true, '2026-01-30 08:44:27.585', '2026-01-30 08:44:27.585', NULL, NULL),
('01KG7BG1Y25P9ZW7MMZ7MGNY7R', 'Bode', true, '2026-01-30 08:44:27.586', '2026-01-30 08:44:27.586', NULL, NULL),
('01KG7BG1Y25P9ZW7MMZ9NHQBZH', 'Ovelha', true, '2026-01-30 08:44:27.586', '2026-01-30 08:44:27.586', NULL, NULL),
('01KG7BG1Y3Y05KSZ6HZQM0R6Z1', 'Vaca', true, '2026-01-30 08:44:27.587', '2026-01-30 08:44:27.587', NULL, NULL),
('01KG7BG1Y3Y05KSZ6HZTBB28P4', 'Búfalo', true, '2026-01-30 08:44:27.587', '2026-01-30 08:44:27.588', NULL, NULL),
('01KG7BG1Y4CB7GPGJVG9839ABK', 'Porco', true, '2026-01-30 08:44:27.588', '2026-01-30 08:44:27.588', NULL, NULL),
('01KG7BG1Y4CB7GPGJVG9T8KCRX', 'Burro', true, '2026-01-30 08:44:27.588', '2026-01-30 08:44:27.589', NULL, NULL),
('01KG7BG1Y5ZAEPVDEKJCRZWEB3', 'Jabuti', true, '2026-01-30 08:44:27.589', '2026-01-30 08:44:27.589', NULL, NULL),
('01KG7BG1Y5ZAEPVDEKJEVZK5M7', 'Caracol', true, '2026-01-30 08:44:27.589', '2026-01-30 08:44:27.589', NULL, NULL),
('01KG7BG1Y5ZAEPVDEKJHMFZ2VC', 'Caranguejo', true, '2026-01-30 08:44:27.589', '2026-01-30 08:44:27.590', NULL, NULL),
('01KG7BG1Y61CAAY6CK1C24YABY', 'Cágado', true, '2026-01-30 08:44:27.590', '2026-01-30 08:44:27.590', NULL, NULL),
('01KG7BG1Y61CAAY6CK1CMS3555', 'Camundongo', true, '2026-01-30 08:44:27.590', '2026-01-30 08:44:27.590', NULL, NULL),
('01KG7BG1Y7ET64T87DZXK3HYPS', 'Gerbil', true, '2026-01-30 08:44:27.591', '2026-01-30 08:44:27.591', NULL, NULL),
('01KG7BG1Y7ET64T87E00YEXN0T', 'Degus', true, '2026-01-30 08:44:27.591', '2026-01-30 08:44:27.591', NULL, NULL),
('01KG7BG1Y7ET64T87E016S91D3', 'Rena', true, '2026-01-30 08:44:27.591', '2026-01-30 08:44:27.591', NULL, NULL),
('01KG7BG1Y8R3S95K9NF3PYE3D3', 'Lhama', true, '2026-01-30 08:44:27.592', '2026-01-30 08:44:27.592', NULL, NULL),
('01KG7BG1Y8R3S95K9NF4HPEDJV', 'Alpaca', true, '2026-01-30 08:44:27.592', '2026-01-30 08:44:27.592', NULL, NULL),
('01KG7BG1Y8R3S95K9NF4XBGCJE', 'Doninha', true, '2026-01-30 08:44:27.592', '2026-01-30 08:44:27.592', NULL, NULL),
('01KG7BG1Y9Y053P2JPNXZ11XS0', 'Morcego', true, '2026-01-30 08:44:27.593', '2026-01-30 08:44:27.593', NULL, NULL),
('01KG7BG1Y9Y053P2JPP0Z7RCZ6', 'Raposa', true, '2026-01-30 08:44:27.593', '2026-01-30 08:44:27.593', NULL, NULL),
('01KG7BG1Y9Y053P2JPP4TS68T7', 'Guaxinim', true, '2026-01-30 08:44:27.593', '2026-01-30 08:44:27.593', NULL, NULL),
('01KG7BG1Y9Y053P2JPP77EV9AJ', 'Esquilo', true, '2026-01-30 08:44:27.593', '2026-01-30 08:44:27.594', NULL, NULL),
('01KG7BG1YA6NKTSGH9SJ78ZJPX', 'Toupeira', true, '2026-01-30 08:44:27.594', '2026-01-30 08:44:27.594', NULL, NULL),
('01KG7BG1YA6NKTSGH9SM1R8R2T', 'Tamanduá', true, '2026-01-30 08:44:27.594', '2026-01-30 08:44:27.594', NULL, NULL),
('01KG7BG1YA6NKTSGH9SPA1V3K7', 'Bicho-preguiça', true, '2026-01-30 08:44:27.594', '2026-01-30 08:44:27.595', NULL, NULL),
('01KG7BG1YBW160041G0VT6XT4N', 'Macaco', true, '2026-01-30 08:44:27.595', '2026-01-30 08:44:27.595', NULL, NULL),
('01KG7BG1YBW160041G0ZFKXK2J', 'Mico', true, '2026-01-30 08:44:27.595', '2026-01-30 08:44:27.595', NULL, NULL),
('01KG7BG1YBW160041G1175PHSH', 'Sagui', true, '2026-01-30 08:44:27.595', '2026-01-30 08:44:27.596', NULL, NULL),
('01KG7BG1YCAR8TGMK2ES2HJ44M', 'Capivara', true, '2026-01-30 08:44:27.596', '2026-01-30 08:44:27.596', NULL, NULL),
('01KG7BG1YCAR8TGMK2ETHEESJD', 'Quati', true, '2026-01-30 08:44:27.596', '2026-01-30 08:44:27.596', NULL, NULL),
('01KG7BG1YCAR8TGMK2EWEYE47T', 'Canguru', true, '2026-01-30 08:44:27.596', '2026-01-30 08:44:27.596', NULL, NULL),
('01KG7BG1YD1K2WR6VK57AJC56V', 'Wallaby', true, '2026-01-30 08:44:27.597', '2026-01-30 08:44:27.597', NULL, NULL),
('01KG7BG1YD1K2WR6VK57WP114D', 'Dromedário', true, '2026-01-30 08:44:27.597', '2026-01-30 08:44:27.597', NULL, NULL),
('01KG7BG1YD1K2WR6VK5BENW692', 'Camelo', true, '2026-01-30 08:44:27.597', '2026-01-30 08:44:27.597', NULL, NULL),
('01KG7BG1YEFGWDPQFCB07VY29H', 'Foca', true, '2026-01-30 08:44:27.598', '2026-01-30 08:44:27.598', NULL, NULL),
('01KG7BG1YEFGWDPQFCB24RK8RG', 'Leão-marinho', true, '2026-01-30 08:44:27.598', '2026-01-30 08:44:27.598', NULL, NULL),
('01KG7BG1YF28TXRCMMX0XQ7164', 'Lontra', true, '2026-01-30 08:44:27.599', '2026-01-30 08:44:27.599', NULL, NULL),
('01KG7BG1YF28TXRCMMX458PW5V', 'Castor', true, '2026-01-30 08:44:27.599', '2026-01-30 08:44:27.599', NULL, NULL),
('01KG7BG1YF28TXRCMMX4EHT9XR', 'Panda', true, '2026-01-30 08:44:27.599', '2026-01-30 08:44:27.600', NULL, NULL),
('01KG7BG1YGE67MAB4PK7MX4M34', 'Urso', true, '2026-01-30 08:44:27.600', '2026-01-30 08:44:27.600', NULL, NULL),
('01KG7BG1YGE67MAB4PKAC402GS', 'Lobo', true, '2026-01-30 08:44:27.600', '2026-01-30 08:44:27.600', NULL, NULL),
('01KG7BG1YGE67MAB4PKAPSEQBH', 'Hiena', true, '2026-01-30 08:44:27.600', '2026-01-30 08:44:27.600', NULL, NULL),
('01KG7BG1YH87N0BTSP1DE7XCV2', 'Chita', true, '2026-01-30 08:44:27.601', '2026-01-30 08:44:27.601', NULL, NULL),
('01KG7BG1YH87N0BTSP1FXQBHA0', 'Leopardo', true, '2026-01-30 08:44:27.601', '2026-01-30 08:44:27.601', NULL, NULL),
('01KG7BG1YH87N0BTSP1KFRCZR7', 'Tigre', true, '2026-01-30 08:44:27.601', '2026-01-30 08:44:27.601', NULL, NULL),
('01KG7BG1YJYAWAVAHZQSH2W30J', 'Leão', true, '2026-01-30 08:44:27.602', '2026-01-30 08:44:27.602', NULL, NULL),
('01KG7BG1YJYAWAVAHZQSVQWH71', 'Guepardo', true, '2026-01-30 08:44:27.602', '2026-01-30 08:44:27.602', NULL, NULL),
('01KG7BG1YJYAWAVAHZQWRHN4QW', 'Pantera', true, '2026-01-30 08:44:27.602', '2026-01-30 08:44:27.602', NULL, NULL),
('01KG7BG1YK02SH92M4QZSK1HP9', 'Onça', true, '2026-01-30 08:44:27.603', '2026-01-30 08:44:27.603', NULL, NULL),
('01KG7BG1YK02SH92M4R278Z8GN', 'Jaguar', true, '2026-01-30 08:44:27.603', '2026-01-30 08:44:27.603', NULL, NULL),
('01KG7BG1YK02SH92M4R5Y8HV94', 'Puma', true, '2026-01-30 08:44:27.603', '2026-01-30 08:44:27.604', NULL, NULL),
('01KG7BG1YMEGJ655M78XA4M560', 'Gato-do-mato', true, '2026-01-30 08:44:27.604', '2026-01-30 08:44:27.604', NULL, NULL),
('01KG7BG1YMEGJ655M78XGK6VET', 'Serval', true, '2026-01-30 08:44:27.604', '2026-01-30 08:44:27.604', NULL, NULL),
('01KG7BG1YNER3PBF6MP0CWRMEK', 'Caracal', true, '2026-01-30 08:44:27.605', '2026-01-30 08:44:27.605', NULL, NULL),
('01KG7BG1YNER3PBF6MP0YB1CN6', 'Lince', true, '2026-01-30 08:44:27.605', '2026-01-30 08:44:27.605', NULL, NULL),
('01KG7BG1YNER3PBF6MP3XT8QC4', 'Gato-selvagem', true, '2026-01-30 08:44:27.605', '2026-01-30 08:44:27.605', NULL, NULL),
('01KG7BG1YPZKNXQ0YRFRR28FYB', 'Coiote', true, '2026-01-30 08:44:27.606', '2026-01-30 08:44:27.606', NULL, NULL),
('01KG7BG1YPZKNXQ0YRFT34GVJG', 'Dingo', true, '2026-01-30 08:44:27.606', '2026-01-30 08:44:27.606', NULL, NULL),
('01KG7BG1YPZKNXQ0YRFX5VMJR4', 'Chacal', true, '2026-01-30 08:44:27.606', '2026-01-30 08:44:27.606', NULL, NULL),
('01KG7BG1YQ0NGT31ZM46RHH44W', 'Cão-selvagem', true, '2026-01-30 08:44:27.607', '2026-01-30 08:44:27.607', NULL, NULL),
('01KG7BG1YQ0NGT31ZM49N13B9B', 'Cão-do-mato', true, '2026-01-30 08:44:27.607', '2026-01-30 08:44:27.607', NULL, NULL),
('01KG7BG1YQ0NGT31ZM4ANP770S', 'Cão-vinagre', true, '2026-01-30 08:44:27.607', '2026-01-30 08:44:27.607', NULL, NULL),
('01KG7BG1YRB1DTKKQV7HAHDH9E', 'Cão-guará', true, '2026-01-30 08:44:27.608', '2026-01-30 08:44:27.608', NULL, NULL),
('01KG7BG1YRB1DTKKQV7JT53QZ1', 'Cão-lobo', true, '2026-01-30 08:44:27.608', '2026-01-30 08:44:27.608', NULL, NULL),
('01KG7BG1YRB1DTKKQV7MVVYTAH', 'Cão-d''água', true, '2026-01-30 08:44:27.608', '2026-01-30 08:44:27.609', NULL, NULL),
('01KG7BG1YSCFHS2WRRN7J07FX9', 'Cão-de-crista', true, '2026-01-30 08:44:27.609', '2026-01-30 08:44:27.609', NULL, NULL),
('01KG7BG1YSCFHS2WRRN7K8A16T', 'Cão-de-montanha', true, '2026-01-30 08:44:27.609', '2026-01-30 08:44:27.609', NULL, NULL),
('01KG7BG1YSCFHS2WRRN8BP3JT8', 'Cão-de-caça', true, '2026-01-30 08:44:27.609', '2026-01-30 08:44:27.610', NULL, NULL),
('01KG7BG1YTE81JJP368E0BYNWE', 'Cão-de-guarda', true, '2026-01-30 08:44:27.610', '2026-01-30 08:44:27.610', NULL, NULL),
('01KG7BG1YTE81JJP368HEDJ7G6', 'Cão-de-pastoreio', true, '2026-01-30 08:44:27.610', '2026-01-30 08:44:27.610', NULL, NULL),
('01KG7BG1YTE81JJP368K1JC0Z2', 'Cão-de-trenó', true, '2026-01-30 08:44:27.610', '2026-01-30 08:44:27.611', NULL, NULL),
('01KG7BG1YVDNNWNRJVNHYEK1SA', 'Cão-de-aporte', true, '2026-01-30 08:44:27.611', '2026-01-30 08:44:27.611', NULL, NULL),
('01KG7BG1YVDNNWNRJVNKAFRH7D', 'Cão-de-pele', true, '2026-01-30 08:44:27.611', '2026-01-30 08:44:27.611', NULL, NULL),
('01KG7BG1YVDNNWNRJVNN7ZKJPZ', 'Cão-de-pelo-curto', true, '2026-01-30 08:44:27.611', '2026-01-30 08:44:27.612', NULL, NULL),
('01KG7BG1YW549XJA4KP33R6R1G', 'Cão-de-pelo-longo', true, '2026-01-30 08:44:27.612', '2026-01-30 08:44:27.612', NULL, NULL),
('01KG7BG1YW549XJA4KP61GE823', 'Cão-de-pelo-ondulado', true, '2026-01-30 08:44:27.612', '2026-01-30 08:44:27.612', NULL, NULL),
('01KG7BG1YW549XJA4KP8668R34', 'Cão-de-pelo-encaracolado', true, '2026-01-30 08:44:27.612', '2026-01-30 08:44:27.612', NULL, NULL),
('01KG7BG1YXRX6RJCRGGK6M85VF', 'Cão-de-pelo-duro', true, '2026-01-30 08:44:27.613', '2026-01-30 08:44:27.613', NULL, NULL),
('01KG7BG1YXRX6RJCRGGN2C4M68', 'Cão-de-pelo-sedoso', true, '2026-01-30 08:44:27.613', '2026-01-30 08:44:27.613', NULL, NULL),
('01KG7BG1YXRX6RJCRGGPC8NN0E', 'Cão-de-pelo-áspero', true, '2026-01-30 08:44:27.613', '2026-01-30 08:44:27.613', NULL, NULL),
('01KG7BG1YYVPWTHRKFQDB34CGJ', 'Cão-de-pelo-macio', true, '2026-01-30 08:44:27.614', '2026-01-30 08:44:27.614', NULL, NULL),
('01KG7BG1YYVPWTHRKFQFG33W81', 'Cão-de-pelo-fino', true, '2026-01-30 08:44:27.614', '2026-01-30 08:44:27.614', NULL, NULL),
('01KG7BG1YYVPWTHRKFQH5F4CS1', 'Cão-de-pelo-grosso', true, '2026-01-30 08:44:27.614', '2026-01-30 08:44:27.614', NULL, NULL),
('01KG7BG1YYVPWTHRKFQHAZH6Z2', 'Cão-de-pelo-espesso', true, '2026-01-30 08:44:27.614', '2026-01-30 08:44:27.615', NULL, NULL),
('01KG7BG1YZT0R78ZJQKV6PGZKG', 'Cão-de-pelo-fofo', true, '2026-01-30 08:44:27.615', '2026-01-30 08:44:27.615', NULL, NULL),
('01KG7BG1YZT0R78ZJQKY5956T6', 'Cão-de-pelo-liso', true, '2026-01-30 08:44:27.615', '2026-01-30 08:44:27.615', NULL, NULL),
('01KG7BG1YZT0R78ZJQKZRMWWDW', 'Cão-de-pelo-ondulado-curto', true, '2026-01-30 08:44:27.615', '2026-01-30 08:44:27.616', NULL, NULL),
('01KG7BG1Z0TS072KYEVM16CBK1', 'Cão-de-pelo-ondulado-longo', true, '2026-01-30 08:44:27.616', '2026-01-30 08:44:27.616', NULL, NULL),
('01KG7BG1Z0TS072KYEVM71F7HW', 'Cão-de-pelo-ondulado-encaracolado', true, '2026-01-30 08:44:27.616', '2026-01-30 08:44:27.616', NULL, NULL),
('01KG7BG1Z0TS072KYEVQEJA9AV', 'Cão-de-pelo-ondulado-duro', true, '2026-01-30 08:44:27.616', '2026-01-30 08:44:27.617', NULL, NULL),
('01KG7BG1Z1K3FPCBM28RZPT9G8', 'Cão-de-pelo-ondulado-sedoso', true, '2026-01-30 08:44:27.617', '2026-01-30 08:44:27.617', NULL, NULL),
('01KG7BG1Z1K3FPCBM28WGETTWA', 'Cão-de-pelo-ondulado-áspero', true, '2026-01-30 08:44:27.617', '2026-01-30 08:44:27.617', NULL, NULL),
('01KG7BG1Z1K3FPCBM28YNBAZBZ', 'Cão-de-pelo-ondulado-macio', true, '2026-01-30 08:44:27.617', '2026-01-30 08:44:27.618', NULL, NULL),
('01KG7BG1Z2S4NNFG5YES16EAG5', 'Cão-de-pelo-ondulado-fino', true, '2026-01-30 08:44:27.618', '2026-01-30 08:44:27.618', NULL, NULL),
('01KG7BG1Z2S4NNFG5YETRY0MY9', 'Cão-de-pelo-ondulado-grosso', true, '2026-01-30 08:44:27.618', '2026-01-30 08:44:27.618', NULL, NULL),
('01KG7BG1Z2S4NNFG5YEV3P13W3', 'Cão-de-pelo-ondulado-espesso', true, '2026-01-30 08:44:27.618', '2026-01-30 08:44:27.619', NULL, NULL),
('01KG7BG1Z3E88MBHKP8K19GRDM', 'Cão-de-pelo-ondulado-fofo', true, '2026-01-30 08:44:27.619', '2026-01-30 08:44:27.619', NULL, NULL),
('01KG7BG1Z3E88MBHKP8N674E40', 'Cão-de-pelo-ondulado-liso', true, '2026-01-30 08:44:27.619', '2026-01-30 08:44:27.619', NULL, NULL)
ON CONFLICT (id) DO NOTHING;
`

	return db.Exec(seed).Error
}

func Migration20260218000001(db *gorm.DB) error {
	const seed = `
INSERT INTO public.categories
(id, "name", active, created_at, updated_at, deactivated_at, deleted_at)
VALUES
('01KG7CG1XN0BQ4KHKPHY5V5CAT', 'Pet Sitter', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N46CPCAT1', 'Passeador', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4ACPCAT2', 'Dog Walker', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4BCPCAT3', 'Pet Shop', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4CCPCAT4', 'Tosa e Banho', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4DCPCAT5', 'Grooming', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4ECPCAT6', 'Veterinário', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4FCPCAT7', 'Clínica Veterinária', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4GCPCAT8', 'Adestramento', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4HCPCAT9', 'Treinamento', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4ICPCA10', 'Hotel para Pets', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4JCPCA11', 'Creche para Pets', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4KCPCA12', 'Day Care', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4LCPCA13', 'Hospedagem', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4MCPCA14', 'Transporte de Pets', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4NCPCA15', 'Pet Taxi', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4OCPCA16', 'Fotografia de Pets', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4PCPCA17', 'Alimentação', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4QCPCA18', 'Nutrição Animal', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4RCPCA19', 'Fisioterapia Veterinária', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4SCPCA20', 'Acupuntura Veterinária', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4TCPCA21', 'Emergência 24h', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4UCPCA22', 'Spa para Pets', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4VCPCA23', 'Hidratação', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4WCPCA24', 'Tosa Higiênica', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4XCPCA25', 'Corte de Unhas', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4YCPCA26', 'Comportamento Animal', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N4ZCPCA27', 'Consulta Veterinária', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N5ACPCA28', 'Vacinação', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N5BCPCA29', 'Cirurgia Veterinária', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N5CCPCA30', 'Castração', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N5DCPCA31', 'Exames Laboratoriais', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N5ECPCA32', 'Ultrassom Veterinário', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N5FCPCA33', 'Raio-X Veterinário', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N5GCPCA34', 'Dentista Veterinário', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL),
('01KG7CG1XR2FA63J6N5HCPCA35', 'Limpeza Dentária', true, '2026-02-18 10:00:00.000', '2026-02-18 10:00:00.000', NULL, NULL)
ON CONFLICT (id) DO NOTHING;
`

	return db.Exec(seed).Error
}

func Migration20260311000000(db *gorm.DB) error {
	ownerPassword, err := bcrypt.GenerateFromPassword([]byte("Owner@123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	providerPassword, err := bcrypt.GenerateFromPassword([]byte("Provider@123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	owner := models.User{
		ID:              "01KPV3M8F8A1B2C3D4E5F6G7H8",
		Name:            "Owner Seed",
		UserType:        "owner",
		Email:           "owner.seed@petservices.local",
		Password:        string(ownerPassword),
		CountryCode:     "+55",
		AreaCode:        "11",
		PhoneNumber:     "999999111",
		EmailVerified:   true,
		ProfileComplete: true,
		Active:          true,
		Street:          "Rua das Flores",
		Number:          "100",
		Neighborhood:    "Centro",
		City:            "Sao Paulo",
		ZipCode:         "01001-000",
		State:           "SP",
		Country:         "BR",
		Complement:      "Apto 12",
		Latitude:        -23.550520,
		Longitude:       -46.633308,
	}

	if err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
			"user_type",
			"password",
			"country_code",
			"area_code",
			"phone_number",
			"email_verified",
			"profile_complete",
			"active",
			"street",
			"number",
			"neighborhood",
			"city",
			"zip_code",
			"state",
			"country",
			"complement",
			"latitude",
			"longitude",
		}),
	}).Create(&owner).Error; err != nil {
		return err
	}

	providerUser := models.User{
		ID:              "01KPV3M8F8J1K2L3M4N5P6Q7R8",
		Name:            "Provider Seed",
		UserType:        "provider",
		Email:           "provider.seed@petservices.local",
		Password:        string(providerPassword),
		CountryCode:     "+55",
		AreaCode:        "11",
		PhoneNumber:     "999999222",
		EmailVerified:   true,
		ProfileComplete: true,
		Active:          true,
		Street:          "Avenida Paulista",
		Number:          "1578",
		Neighborhood:    "Bela Vista",
		City:            "Sao Paulo",
		ZipCode:         "01310-200",
		State:           "SP",
		Country:         "BR",
		Complement:      "Sala 305",
		Latitude:        -23.561400,
		Longitude:       -46.655881,
	}

	if err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
			"user_type",
			"password",
			"country_code",
			"area_code",
			"phone_number",
			"email_verified",
			"profile_complete",
			"active",
			"street",
			"number",
			"neighborhood",
			"city",
			"zip_code",
			"state",
			"country",
			"complement",
			"latitude",
			"longitude",
		}),
	}).Create(&providerUser).Error; err != nil {
		return err
	}

	var persistedProviderUser models.User
	if err := db.Where("email = ?", providerUser.Email).First(&persistedProviderUser).Error; err != nil {
		return err
	}

	provider := models.Provider{
		ID:            "01KPV3M8F8S1T2U3V4W5X6Y7Z8",
		UserID:        persistedProviderUser.ID,
		BusinessName:  "Provider Seed Pet Care",
		Description:   "Perfil seed de provider para desenvolvimento local",
		PriceRange:    "medium",
		AverageRating: 0.0,
		Active:        true,
		Street:        "Avenida Paulista",
		Number:        "1578",
		Neighborhood:  "Bela Vista",
		City:          "Sao Paulo",
		ZipCode:       "01310-200",
		State:         "SP",
		Country:       "BR",
		Complement:    "Sala 305",
		Latitude:      -23.561400,
		Longitude:     -46.655881,
	}

	if err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"business_name",
			"description",
			"price_range",
			"average_rating",
			"active",
			"street",
			"number",
			"neighborhood",
			"city",
			"zip_code",
			"state",
			"country",
			"complement",
			"latitude",
			"longitude",
		}),
	}).Create(&provider).Error; err != nil {
		return err
	}

	return nil
}

func Migration20260311000001(db *gorm.DB) error {
	var ownerUser models.User
	if err := db.Where("email = ?", "owner.seed@petservices.local").First(&ownerUser).Error; err != nil {
		return err
	}

	pets := []models.Pet{
		{
			ID:        "01KPV6J8A1B2C3D4E5F6G7H8J9",
			UserID:    ownerUser.ID,
			Name:      "Thor",
			SpeciesID: "01KG7BG1XN0BQ4KHKPHY5V5ZEW", // Cachorro
			Breed:     "Labrador",
			Age:       3,
			Weight:    18.40,
			Notes:     "Pet seed owner - cachorro",
			Active:    true,
		},
		{
			ID:        "01KPV6J8A1K2L3M4N5P6Q7R8S9",
			UserID:    ownerUser.ID,
			Name:      "Mia",
			SpeciesID: "01KG7BG1XR2FA63J6N46CPVCKS", // Gato
			Breed:     "Siamês",
			Age:       2,
			Weight:    4.70,
			Notes:     "Pet seed owner - gato",
			Active:    true,
		},
	}

	for _, pet := range pets {
		if err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"user_id",
				"name",
				"species_id",
				"breed",
				"age",
				"weight",
				"notes",
				"active",
			}),
		}).Create(&pet).Error; err != nil {
			return err
		}
	}

	return nil
}

func Migration20260311000002(db *gorm.DB) error {
	seedTags := []models.Tag{
		{ID: "01KTAG1XN0BQ4KHKPHY5VTAG01", Name: "Domiciliar", Active: true},
		{ID: "01KTAG1XN0BQ4KHKPHY5VTAG02", Name: "Pequeno Porte", Active: true},
		{ID: "01KTAG1XN0BQ4KHKPHY5VTAG03", Name: "Grande Porte", Active: true},
		{ID: "01KTAG1XN0BQ4KHKPHY5VTAG04", Name: "Emergência", Active: true},
		{ID: "01KTAG1XN0BQ4KHKPHY5VTAG05", Name: "Fins de Semana", Active: true},
	}

	for _, tag := range seedTags {
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		}).Create(&tag).Error; err != nil {
			return err
		}
	}

	var providerUser models.User
	if err := db.Where("email = ?", "provider.seed@petservices.local").First(&providerUser).Error; err != nil {
		return err
	}

	var provider models.Provider
	if err := db.Where("user_id = ?", providerUser.ID).First(&provider).Error; err != nil {
		return err
	}

	type serviceSeed struct {
		ID           string
		Name         string
		Description  string
		Price        float64
		PriceMinimum float64
		PriceMaximum float64
		Duration     int
		CategoryIDs  []string
		TagIDs       []string
	}

	serviceSeeds := []serviceSeed{
		{
			ID:           "01KSVC1XN0BQ4KHKPHY5VSVC01",
			Name:         "Banho e Tosa Completo",
			Description:  "Serviço completo de banho e tosa para cães e gatos. Inclui secagem, escovação e perfume.",
			Price:        80.00,
			PriceMinimum: 60.00,
			PriceMaximum: 120.00,
			Duration:     90,
			CategoryIDs:  []string{"01KG7CG1XR2FA63J6N4CCPCAT4", "01KG7CG1XR2FA63J6N4DCPCAT5"},
			TagIDs:       []string{"01KTAG1XN0BQ4KHKPHY5VTAG02", "01KTAG1XN0BQ4KHKPHY5VTAG03"},
		},
		{
			ID:           "01KSVC1XN0BQ4KHKPHY5VSVC02",
			Name:         "Passeio Diário",
			Description:  "Passeio diário de 30 a 60 minutos em parques e áreas verdes próximas.",
			Price:        40.00,
			PriceMinimum: 35.00,
			PriceMaximum: 55.00,
			Duration:     45,
			CategoryIDs:  []string{"01KG7CG1XR2FA63J6N46CPCAT1", "01KG7CG1XR2FA63J6N4ACPCAT2"},
			TagIDs:       []string{"01KTAG1XN0BQ4KHKPHY5VTAG02", "01KTAG1XN0BQ4KHKPHY5VTAG05"},
		},
		{
			ID:           "01KSVC1XN0BQ4KHKPHY5VSVC03",
			Name:         "Pet Sitting em Casa",
			Description:  "Cuidado do seu pet no conforto da sua própria casa. Inclui alimentação, brincadeiras e companhia.",
			Price:        120.00,
			PriceMinimum: 100.00,
			PriceMaximum: 180.00,
			Duration:     480,
			CategoryIDs:  []string{"01KG7CG1XN0BQ4KHKPHY5V5CAT"},
			TagIDs:       []string{"01KTAG1XN0BQ4KHKPHY5VTAG01", "01KTAG1XN0BQ4KHKPHY5VTAG05"},
		},
		{
			ID:           "01KSVC1XN0BQ4KHKPHY5VSVC04",
			Name:         "Hospedagem (Hotel Pet)",
			Description:  "Hospedagem completa com alimentação, passeios e muita atenção. Ideal para viagens.",
			Price:        150.00,
			PriceMinimum: 120.00,
			PriceMaximum: 200.00,
			Duration:     1440,
			CategoryIDs:  []string{"01KG7CG1XR2FA63J6N4ICPCA10"},
			TagIDs:       []string{"01KTAG1XN0BQ4KHKPHY5VTAG02", "01KTAG1XN0BQ4KHKPHY5VTAG03", "01KTAG1XN0BQ4KHKPHY5VTAG05"},
		},
		{
			ID:           "01KSVC1XN0BQ4KHKPHY5VSVC05",
			Name:         "Adestramento Básico",
			Description:  "Treinamento de comandos básicos: senta, fica, deita, não pula. Sessões de 1 hora com reforço positivo.",
			Price:        200.00,
			PriceMinimum: 180.00,
			PriceMaximum: 250.00,
			Duration:     60,
			CategoryIDs:  []string{"01KG7CG1XR2FA63J6N4GCPCAT8", "01KG7CG1XR2FA63J6N4HCPCAT9"},
			TagIDs:       []string{"01KTAG1XN0BQ4KHKPHY5VTAG02", "01KTAG1XN0BQ4KHKPHY5VTAG03"},
		},
	}

	for _, svc := range serviceSeeds {
		categories := make([]models.Category, len(svc.CategoryIDs))
		for i, id := range svc.CategoryIDs {
			categories[i] = models.Category{ID: id}
		}

		tagModels := make([]models.Tag, len(svc.TagIDs))
		for i, id := range svc.TagIDs {
			tagModels[i] = models.Tag{ID: id}
		}

		service := models.Service{
			ID:           svc.ID,
			ProviderID:   provider.ID,
			Name:         svc.Name,
			Description:  svc.Description,
			Price:        svc.Price,
			PriceMinimum: svc.PriceMinimum,
			PriceMaximum: svc.PriceMaximum,
			Duration:     svc.Duration,
			Active:       true,
		}

		if err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"provider_id", "name", "description", "price",
				"price_minimum", "price_maximum", "duration", "active",
			}),
		}).Create(&service).Error; err != nil {
			return err
		}

		if err := db.Model(&service).Association("Categories").Replace(categories); err != nil {
			return err
		}

		if err := db.Model(&service).Association("Tags").Replace(tagModels); err != nil {
			return err
		}
	}

	return nil
}
func Migration20260315000000(db *gorm.DB) error {
	providerPassword, err := bcrypt.GenerateFromPassword([]byte("Provider@123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	type svcSeed struct {
		ID           string
		Name         string
		Description  string
		Price        float64
		PriceMinimum float64
		PriceMaximum float64
		Duration     int
		CategoryIDs  []string
		TagIDs       []string
	}

	type seedEntry struct {
		User     models.User
		Provider models.Provider
		Services []svcSeed
	}

	pwd := string(providerPassword)

	entries := []seedEntry{
		// 1. Ana Beatriz Santos — Pet Sitter / Passeio — Centro, Aracaju
		{
			User: models.User{
				ID: "01KSSEPROVARACAJU00USER001", Name: "Ana Beatriz Santos",
				UserType: "provider", Email: "ana.beatriz.provider@petservices.local",
				Password: pwd, CountryCode: "+55", AreaCode: "79", PhoneNumber: "991110001",
				EmailVerified: true, ProfileComplete: true, Active: true,
				Street: "Rua São João", Number: "245", Neighborhood: "Centro",
				City: "Aracaju", ZipCode: "49010-020", State: "SE", Country: "BR",
				Latitude: -10.9099, Longitude: -37.0504,
			},
			Provider: models.Provider{
				ID: "01KSSEPROVARACAJU00PROV001", BusinessName: "Ana Pet Care",
				Description:   "Serviços de pet sitting e passeios em Aracaju. Cuidado especial para cães e gatos no Centro da cidade.",
				PriceRange:    "low",
				AverageRating: 0.0, Active: true,
				Street: "Rua São João", Number: "245", Neighborhood: "Centro",
				City: "Aracaju", ZipCode: "49010-020", State: "SE", Country: "BR",
				Latitude: -10.9099, Longitude: -37.0504,
			},
			Services: []svcSeed{
				{
					ID:          "01KSSEPROVARACAJU00SVC0101",
					Name:        "Pet Sitting em Casa",
					Description: "Cuidado do seu pet no conforto da sua própria casa. Alimentação, brincadeiras e atenção o dia todo.",
					Price:       100.00, PriceMinimum: 80.00, PriceMaximum: 150.00, Duration: 480,
					CategoryIDs: []string{"01KG7CG1XN0BQ4KHKPHY5V5CAT"},
					TagIDs:      []string{"01KTAG1XN0BQ4KHKPHY5VTAG01", "01KTAG1XN0BQ4KHKPHY5VTAG05"},
				},
				{
					ID:          "01KSSEPROVARACAJU00SVC0102",
					Name:        "Passeio Diário",
					Description: "Passeios diários pelo Parque da Sementeira e orla de Aracaju. Duração de 45 a 60 minutos.",
					Price:       35.00, PriceMinimum: 30.00, PriceMaximum: 50.00, Duration: 45,
					CategoryIDs: []string{"01KG7CG1XR2FA63J6N46CPCAT1", "01KG7CG1XR2FA63J6N4ACPCAT2"},
					TagIDs:      []string{"01KTAG1XN0BQ4KHKPHY5VTAG02", "01KTAG1XN0BQ4KHKPHY5VTAG05"},
				},
			},
		},
		// 2. Carlos Eduardo Oliveira — Veterinário — Farolândia, Aracaju
		{
			User: models.User{
				ID: "01KSSEPROVARACAJU00USER002", Name: "Carlos Eduardo Oliveira",
				UserType: "provider", Email: "carlos.eduardo.provider@petservices.local",
				Password: pwd, CountryCode: "+55", AreaCode: "79", PhoneNumber: "991110002",
				EmailVerified: true, ProfileComplete: true, Active: true,
				Street: "Avenida Beira Mar", Number: "850", Neighborhood: "Farolândia",
				City: "Aracaju", ZipCode: "49030-100", State: "SE", Country: "BR",
				Latitude: -10.9581, Longitude: -37.0568,
			},
			Provider: models.Provider{
				ID: "01KSSEPROVARACAJU00PROV002", BusinessName: "Clínica Vet Farolândia",
				Description:   "Clínica veterinária completa no bairro Farolândia. Especializada em pequenos animais com atendimento humanizado.",
				PriceRange:    "high",
				AverageRating: 0.0, Active: true,
				Street: "Avenida Beira Mar", Number: "850", Neighborhood: "Farolândia",
				City: "Aracaju", ZipCode: "49030-100", State: "SE", Country: "BR",
				Latitude: -10.9581, Longitude: -37.0568,
			},
			Services: []svcSeed{
				{
					ID:          "01KSSEPROVARACAJU00SVC0201",
					Name:        "Consulta Veterinária",
					Description: "Consulta clínica geral com anamnese completa, exame físico e orientações de saúde para seu pet.",
					Price:       180.00, PriceMinimum: 150.00, PriceMaximum: 220.00, Duration: 45,
					CategoryIDs: []string{"01KG7CG1XR2FA63J6N4ECPCAT6", "01KG7CG1XR2FA63J6N4FCPCAT7"},
					TagIDs:      []string{"01KTAG1XN0BQ4KHKPHY5VTAG02", "01KTAG1XN0BQ4KHKPHY5VTAG03"},
				},
				{
					ID:          "01KSSEPROVARACAJU00SVC0202",
					Name:        "Vacinação",
					Description: "Aplicação de vacinas essenciais e controle do cartão de vacinação do seu animal.",
					Price:       120.00, PriceMinimum: 90.00, PriceMaximum: 160.00, Duration: 30,
					CategoryIDs: []string{"01KG7CG1XR2FA63J6N4ECPCAT6", "01KG7CG1XR2FA63J6N5ACPCA28"},
					TagIDs:      []string{"01KTAG1XN0BQ4KHKPHY5VTAG02", "01KTAG1XN0BQ4KHKPHY5VTAG03"},
				},
				{
					ID:          "01KSSEPROVARACAJU00SVC0203",
					Name:        "Emergência 24h",
					Description: "Atendimento de emergência disponível 24 horas para situações críticas que exigem cuidados imediatos.",
					Price:       350.00, PriceMinimum: 250.00, PriceMaximum: 600.00, Duration: 60,
					CategoryIDs: []string{"01KG7CG1XR2FA63J6N4ECPCAT6", "01KG7CG1XR2FA63J6N4TCPCA21"},
					TagIDs:      []string{"01KTAG1XN0BQ4KHKPHY5VTAG04"},
				},
			},
		},
		// 3. Mariana Costa — Banho e Tosa / Grooming — Atalaia, Aracaju
		{
			User: models.User{
				ID: "01KSSEPROVARACAJU00USER003", Name: "Mariana Costa",
				UserType: "provider", Email: "mariana.costa.provider@petservices.local",
				Password: pwd, CountryCode: "+55", AreaCode: "79", PhoneNumber: "991110003",
				EmailVerified: true, ProfileComplete: true, Active: true,
				Street: "Avenida Oceanica", Number: "312", Neighborhood: "Atalaia",
				City: "Aracaju", ZipCode: "49037-100", State: "SE", Country: "BR",
				Latitude: -10.9856, Longitude: -37.0455,
			},
			Provider: models.Provider{
				ID: "01KSSEPROVARACAJU00PROV003", BusinessName: "Mariana Grooming Atalaia",
				Description:   "Pet shop e grooming especializado na Atalaia. Banho, tosa e estética para seu pet com produtos premium.",
				PriceRange:    "medium",
				AverageRating: 0.0, Active: true,
				Street: "Avenida Oceanica", Number: "312", Neighborhood: "Atalaia",
				City: "Aracaju", ZipCode: "49037-100", State: "SE", Country: "BR",
				Latitude: -10.9856, Longitude: -37.0455,
			},
			Services: []svcSeed{
				{
					ID:          "01KSSEPROVARACAJU00SVC0301",
					Name:        "Banho e Tosa Completo",
					Description: "Serviço completo de banho e tosa para cães e gatos. Inclui secagem, escovação, limpeza de ouvidos e perfume.",
					Price:       85.00, PriceMinimum: 65.00, PriceMaximum: 130.00, Duration: 90,
					CategoryIDs: []string{"01KG7CG1XR2FA63J6N4CCPCAT4", "01KG7CG1XR2FA63J6N4DCPCAT5"},
					TagIDs:      []string{"01KTAG1XN0BQ4KHKPHY5VTAG02", "01KTAG1XN0BQ4KHKPHY5VTAG03"},
				},
				{
					ID:          "01KSSEPROVARACAJU00SVC0302",
					Name:        "Tosa Higiênica",
					Description: "Tosa higiênica nas regiões de patas, virilha, barriga e ao redor do ânus. Mantém o pet limpo e confortável.",
					Price:       55.00, PriceMinimum: 40.00, PriceMaximum: 75.00, Duration: 45,
					CategoryIDs: []string{"01KG7CG1XR2FA63J6N4CCPCAT4", "01KG7CG1XR2FA63J6N4WCPCA24"},
					TagIDs:      []string{"01KTAG1XN0BQ4KHKPHY5VTAG02"},
				},
				{
					ID:          "01KSSEPROVARACAJU00SVC0303",
					Name:        "Corte de Unhas",
					Description: "Corte e lixamento de unhas para cães e gatos. Procedimento rápido e seguro.",
					Price:       25.00, PriceMinimum: 20.00, PriceMaximum: 35.00, Duration: 20,
					CategoryIDs: []string{"01KG7CG1XR2FA63J6N4DCPCAT5", "01KG7CG1XR2FA63J6N4XCPCA25"},
					TagIDs:      []string{"01KTAG1XN0BQ4KHKPHY5VTAG02", "01KTAG1XN0BQ4KHKPHY5VTAG03"},
				},
			},
		},
		// 4. Rafael Mendes — Adestramento — Jabotiana, Aracaju
		{
			User: models.User{
				ID: "01KSSEPROVARACAJU00USER004", Name: "Rafael Mendes",
				UserType: "provider", Email: "rafael.mendes.provider@petservices.local",
				Password: pwd, CountryCode: "+55", AreaCode: "79", PhoneNumber: "991110004",
				EmailVerified: true, ProfileComplete: true, Active: true,
				Street: "Rua Prefeito Olavo Bilac", Number: "78", Neighborhood: "Jabotiana",
				City: "Aracaju", ZipCode: "49095-000", State: "SE", Country: "BR",
				Latitude: -10.9311, Longitude: -37.0820,
			},
			Provider: models.Provider{
				ID: "01KSSEPROVARACAJU00PROV004", BusinessName: "Rafael Adestra Pets",
				Description:   "Adestramento e treinamento comportamental para cães de todos os portes em Aracaju. Método de reforço positivo.",
				PriceRange:    "medium",
				AverageRating: 0.0, Active: true,
				Street: "Rua Prefeito Olavo Bilac", Number: "78", Neighborhood: "Jabotiana",
				City: "Aracaju", ZipCode: "49095-000", State: "SE", Country: "BR",
				Latitude: -10.9311, Longitude: -37.0820,
			},
			Services: []svcSeed{
				{
					ID:          "01KSSEPROVARACAJU00SVC0401",
					Name:        "Adestramento Básico",
					Description: "Treino de comandos básicos: senta, fica, deita, não pula. Sessões de 1 hora com reforço positivo.",
					Price:       200.00, PriceMinimum: 180.00, PriceMaximum: 250.00, Duration: 60,
					CategoryIDs: []string{"01KG7CG1XR2FA63J6N4GCPCAT8", "01KG7CG1XR2FA63J6N4HCPCAT9"},
					TagIDs:      []string{"01KTAG1XN0BQ4KHKPHY5VTAG02", "01KTAG1XN0BQ4KHKPHY5VTAG03"},
				},
				{
					ID:          "01KSSEPROVARACAJU00SVC0402",
					Name:        "Treinamento Avançado",
					Description: "Treinamento avançado com foco em comportamento, socialização e obediência em ambientes externos.",
					Price:       280.00, PriceMinimum: 250.00, PriceMaximum: 350.00, Duration: 90,
					CategoryIDs: []string{"01KG7CG1XR2FA63J6N4GCPCAT8", "01KG7CG1XR2FA63J6N4HCPCAT9", "01KG7CG1XR2FA63J6N4YCPCA26"},
					TagIDs:      []string{"01KTAG1XN0BQ4KHKPHY5VTAG02", "01KTAG1XN0BQ4KHKPHY5VTAG03"},
				},
			},
		},
		// 5. Fernanda Lima — Hotel para Pets / Day Care — Coroa do Meio, Aracaju
		{
			User: models.User{
				ID: "01KSSEPROVARACAJU00USER005", Name: "Fernanda Lima",
				UserType: "provider", Email: "fernanda.lima.provider@petservices.local",
				Password: pwd, CountryCode: "+55", AreaCode: "79", PhoneNumber: "991110005",
				EmailVerified: true, ProfileComplete: true, Active: true,
				Street: "Rua Niceu Dantas", Number: "150", Neighborhood: "Coroa do Meio",
				City: "Aracaju", ZipCode: "49025-000", State: "SE", Country: "BR",
				Latitude: -10.9681, Longitude: -37.0336,
			},
			Provider: models.Provider{
				ID: "01KSSEPROVARACAJU00PROV005", BusinessName: "Hotel Pets Coroa do Meio",
				Description:   "Hotel e creche para pets na orla de Aracaju. Hospedagem completa com alimentação individualizada e muita atenção.",
				PriceRange:    "medium",
				AverageRating: 0.0, Active: true,
				Street: "Rua Niceu Dantas", Number: "150", Neighborhood: "Coroa do Meio",
				City: "Aracaju", ZipCode: "49025-000", State: "SE", Country: "BR",
				Latitude: -10.9681, Longitude: -37.0336,
			},
			Services: []svcSeed{
				{
					ID:          "01KSSEPROVARACAJU00SVC0501",
					Name:        "Hospedagem (Hotel Pet)",
					Description: "Hospedagem completa com alimentação, passeios diários e muita atenção. Ideal para viagens e fins de semana.",
					Price:       160.00, PriceMinimum: 130.00, PriceMaximum: 210.00, Duration: 1440,
					CategoryIDs: []string{"01KG7CG1XR2FA63J6N4ICPCA10", "01KG7CG1XR2FA63J6N4LCPCA13"},
					TagIDs:      []string{"01KTAG1XN0BQ4KHKPHY5VTAG02", "01KTAG1XN0BQ4KHKPHY5VTAG03", "01KTAG1XN0BQ4KHKPHY5VTAG05"},
				},
				{
					ID:          "01KSSEPROVARACAJU00SVC0502",
					Name:        "Day Care para Pets",
					Description: "Creche diurna com atividades lúdicas, socialização e alimentação. O pet fica feliz enquanto você trabalha.",
					Price:       80.00, PriceMinimum: 65.00, PriceMaximum: 100.00, Duration: 480,
					CategoryIDs: []string{"01KG7CG1XR2FA63J6N4JCPCA11", "01KG7CG1XR2FA63J6N4KCPCA12"},
					TagIDs:      []string{"01KTAG1XN0BQ4KHKPHY5VTAG02", "01KTAG1XN0BQ4KHKPHY5VTAG03"},
				},
				{
					ID:          "01KSSEPROVARACAJU00SVC0503",
					Name:        "Pet Taxi",
					Description: "Transporte seguro do seu pet para consultas, banhos ou hospedagem. Veículo adaptado e motorista experiente.",
					Price:       60.00, PriceMinimum: 45.00, PriceMaximum: 90.00, Duration: 30,
					CategoryIDs: []string{"01KG7CG1XR2FA63J6N4MCPCA14", "01KG7CG1XR2FA63J6N4NCPCA15"},
					TagIDs:      []string{"01KTAG1XN0BQ4KHKPHY5VTAG02", "01KTAG1XN0BQ4KHKPHY5VTAG03"},
				},
			},
		},
	}

	for _, entry := range entries {
		user := entry.User
		if err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "email"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name", "user_type", "password", "country_code", "area_code", "phone_number",
				"email_verified", "profile_complete", "active", "street", "number",
				"neighborhood", "city", "zip_code", "state", "country", "complement",
				"latitude", "longitude",
			}),
		}).Create(&user).Error; err != nil {
			return err
		}

		var persistedUser models.User
		if err := db.Where("email = ?", user.Email).First(&persistedUser).Error; err != nil {
			return err
		}

		provider := entry.Provider
		provider.UserID = persistedUser.ID
		if err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"business_name", "description", "price_range", "average_rating", "active",
				"street", "number", "neighborhood", "city", "zip_code", "state",
				"country", "complement", "latitude", "longitude",
			}),
		}).Create(&provider).Error; err != nil {
			return err
		}

		var persistedProvider models.Provider
		if err := db.Where("user_id = ?", persistedUser.ID).First(&persistedProvider).Error; err != nil {
			return err
		}

		for _, svc := range entry.Services {
			categories := make([]models.Category, len(svc.CategoryIDs))
			for i, id := range svc.CategoryIDs {
				categories[i] = models.Category{ID: id}
			}
			tagModels := make([]models.Tag, len(svc.TagIDs))
			for i, id := range svc.TagIDs {
				tagModels[i] = models.Tag{ID: id}
			}
			service := models.Service{
				ID:           svc.ID,
				ProviderID:   persistedProvider.ID,
				Name:         svc.Name,
				Description:  svc.Description,
				Price:        svc.Price,
				PriceMinimum: svc.PriceMinimum,
				PriceMaximum: svc.PriceMaximum,
				Duration:     svc.Duration,
				Active:       true,
			}
			if err := db.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"provider_id", "name", "description", "price",
					"price_minimum", "price_maximum", "duration", "active",
				}),
			}).Create(&service).Error; err != nil {
				return err
			}
			if err := db.Model(&service).Association("Categories").Replace(categories); err != nil {
				return err
			}
			if err := db.Model(&service).Association("Tags").Replace(tagModels); err != nil {
				return err
			}
		}
	}

	return nil
}

func Migration20260315000001(db *gorm.DB) error {
	ownerPassword, err := bcrypt.GenerateFromPassword([]byte("Owner@123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	type ownerSeed struct {
		User models.User
	}

	ownerEntries := []ownerSeed{
		{
			User: models.User{
				ID:              "01KOWNU100000000000000001",
				Name:            "Juliana Ferreira",
				UserType:        "owner",
				Email:           "juliana.ferreira.owner@petservices.local",
				Password:        string(ownerPassword),
				CountryCode:     "+55",
				AreaCode:        "79",
				PhoneNumber:     "991120101",
				EmailVerified:   true,
				ProfileComplete: true,
				Active:          true,
				Street:          "Rua Itabaiana",
				Number:          "144",
				Neighborhood:    "Centro",
				City:            "Aracaju",
				ZipCode:         "49010-170",
				State:           "SE",
				Country:         "BR",
				Complement:      "Casa",
				Latitude:        -10.9119,
				Longitude:       -37.0519,
			},
		},
		{
			User: models.User{
				ID:              "01KOWNU100000000000000002",
				Name:            "Paulo Henrique Dias",
				UserType:        "owner",
				Email:           "paulo.dias.owner@petservices.local",
				Password:        string(ownerPassword),
				CountryCode:     "+55",
				AreaCode:        "79",
				PhoneNumber:     "991120102",
				EmailVerified:   true,
				ProfileComplete: true,
				Active:          true,
				Street:          "Avenida Hermes Fontes",
				Number:          "930",
				Neighborhood:    "Salgado Filho",
				City:            "Aracaju",
				ZipCode:         "49020-550",
				State:           "SE",
				Country:         "BR",
				Complement:      "Apto 402",
				Latitude:        -10.9325,
				Longitude:       -37.0586,
			},
		},
		{
			User: models.User{
				ID:              "01KOWNU100000000000000003",
				Name:            "Camila Nunes",
				UserType:        "owner",
				Email:           "camila.nunes.owner@petservices.local",
				Password:        string(ownerPassword),
				CountryCode:     "+55",
				AreaCode:        "79",
				PhoneNumber:     "991120103",
				EmailVerified:   true,
				ProfileComplete: true,
				Active:          true,
				Street:          "Rua Delmiro Gouveia",
				Number:          "81",
				Neighborhood:    "Atalaia",
				City:            "Aracaju",
				ZipCode:         "49037-530",
				State:           "SE",
				Country:         "BR",
				Complement:      "Casa 2",
				Latitude:        -10.9861,
				Longitude:       -37.0484,
			},
		},
		{
			User: models.User{
				ID:              "01KOWNU100000000000000004",
				Name:            "Diego Matos",
				UserType:        "owner",
				Email:           "diego.matos.owner@petservices.local",
				Password:        string(ownerPassword),
				CountryCode:     "+55",
				AreaCode:        "79",
				PhoneNumber:     "991120104",
				EmailVerified:   true,
				ProfileComplete: true,
				Active:          true,
				Street:          "Rua Porto da Folha",
				Number:          "302",
				Neighborhood:    "Siqueira Campos",
				City:            "Aracaju",
				ZipCode:         "49075-070",
				State:           "SE",
				Country:         "BR",
				Complement:      "Fundos",
				Latitude:        -10.9184,
				Longitude:       -37.0738,
			},
		},
	}

	ownersByEmail := make(map[string]models.User, len(ownerEntries))
	for _, entry := range ownerEntries {
		user := entry.User
		if err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "email"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name", "user_type", "password", "country_code", "area_code", "phone_number",
				"email_verified", "profile_complete", "active", "street", "number",
				"neighborhood", "city", "zip_code", "state", "country", "complement",
				"latitude", "longitude",
			}),
		}).Create(&user).Error; err != nil {
			return err
		}

		var persistedUser models.User
		if err := db.Where("email = ?", user.Email).First(&persistedUser).Error; err != nil {
			return err
		}
		ownersByEmail[user.Email] = persistedUser
	}

	type petSeed struct {
		ID         string
		OwnerEmail string
		Name       string
		SpeciesID  string
		Breed      string
		Age        int
		Weight     float64
		Notes      string
	}

	petSeeds := []petSeed{
		{
			ID:         "01KOWNP100000000000000001",
			OwnerEmail: "juliana.ferreira.owner@petservices.local",
			Name:       "Bento",
			SpeciesID:  "01KG7BG1XN0BQ4KHKPHY5V5ZEW",
			Breed:      "Beagle",
			Age:        4,
			Weight:     12.30,
			Notes:      "Cachorro dócil, se adapta bem a passeios.",
		},
		{
			ID:         "01KOWNP100000000000000002",
			OwnerEmail: "juliana.ferreira.owner@petservices.local",
			Name:       "Nina",
			SpeciesID:  "01KG7BG1XR2FA63J6N46CPVCKS",
			Breed:      "Persa",
			Age:        2,
			Weight:     4.10,
			Notes:      "Gata tranquila, gosta de rotina estável.",
		},
		{
			ID:         "01KOWNP100000000000000003",
			OwnerEmail: "paulo.dias.owner@petservices.local",
			Name:       "Thor",
			SpeciesID:  "01KG7BG1XN0BQ4KHKPHY5V5ZEW",
			Breed:      "Pastor Alemão",
			Age:        5,
			Weight:     21.80,
			Notes:      "Cão ativo, precisa de enriquecimento diário.",
		},
		{
			ID:         "01KOWNP100000000000000004",
			OwnerEmail: "camila.nunes.owner@petservices.local",
			Name:       "Lua",
			SpeciesID:  "01KG7BG1XR2FA63J6N46CPVCKS",
			Breed:      "Maine Coon",
			Age:        3,
			Weight:     4.90,
			Notes:      "Gata curiosa e sociável.",
		},
		{
			ID:         "01KOWNP100000000000000005",
			OwnerEmail: "camila.nunes.owner@petservices.local",
			Name:       "Bob",
			SpeciesID:  "01KG7BG1XN0BQ4KHKPHY5V5ZEW",
			Breed:      "Golden Retriever",
			Age:        1,
			Weight:     9.70,
			Notes:      "Filhote, em fase de socialização.",
		},
		{
			ID:         "01KOWNP100000000000000006",
			OwnerEmail: "diego.matos.owner@petservices.local",
			Name:       "Mel",
			SpeciesID:  "01KG7BG1XN0BQ4KHKPHY5V5ZEW",
			Breed:      "Vira-lata",
			Age:        6,
			Weight:     14.20,
			Notes:      "Cadela calma, responde bem a comandos básicos.",
		},
	}

	petsByID := make(map[string]models.Pet, len(petSeeds))
	for _, seed := range petSeeds {
		owner, ok := ownersByEmail[seed.OwnerEmail]
		if !ok {
			return gorm.ErrRecordNotFound
		}

		pet := models.Pet{
			ID:        seed.ID,
			UserID:    owner.ID,
			Name:      seed.Name,
			SpeciesID: seed.SpeciesID,
			Breed:     seed.Breed,
			Age:       seed.Age,
			Weight:    seed.Weight,
			Notes:     seed.Notes,
			Active:    true,
		}

		if err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"user_id", "name", "species_id", "breed", "age", "weight", "notes", "active",
			}),
		}).Create(&pet).Error; err != nil {
			return err
		}

		var persistedPet models.Pet
		if err := db.Where("id = ?", pet.ID).First(&persistedPet).Error; err != nil {
			return err
		}
		petsByID[pet.ID] = persistedPet
	}

	providerEmails := []string{
		"ana.beatriz.provider@petservices.local",
		"carlos.eduardo.provider@petservices.local",
		"mariana.costa.provider@petservices.local",
		"rafael.mendes.provider@petservices.local",
		"fernanda.lima.provider@petservices.local",
	}

	providersByEmail := make(map[string]models.Provider, len(providerEmails))
	for _, email := range providerEmails {
		var providerUser models.User
		if err := db.Where("email = ?", email).First(&providerUser).Error; err != nil {
			return err
		}

		var provider models.Provider
		if err := db.Where("user_id = ?", providerUser.ID).First(&provider).Error; err != nil {
			return err
		}

		providersByEmail[email] = provider
	}

	type requestSeed struct {
		ID            string
		OwnerEmail    string
		ProviderEmail string
		ServiceName   string
		PetID         string
		Notes         string
		Status        string
		RejectReason  string
	}

	requestSeeds := []requestSeed{
		{
			ID:            "01KOWNR100000000000000001",
			OwnerEmail:    "juliana.ferreira.owner@petservices.local",
			ProviderEmail: "ana.beatriz.provider@petservices.local",
			ServiceName:   "Passeio Diário",
			PetID:         "01KOWNP100000000000000001",
			Notes:         "Preciso de passeios de segunda a sexta pela manhã.",
			Status:        "completed",
			RejectReason:  "",
		},
		{
			ID:            "01KOWNR100000000000000002",
			OwnerEmail:    "juliana.ferreira.owner@petservices.local",
			ProviderEmail: "fernanda.lima.provider@petservices.local",
			ServiceName:   "Day Care para Pets",
			PetID:         "01KOWNP100000000000000002",
			Notes:         "Gostaria de iniciar com dois dias por semana.",
			Status:        "pending",
			RejectReason:  "",
		},
		{
			ID:            "01KOWNR100000000000000003",
			OwnerEmail:    "paulo.dias.owner@petservices.local",
			ProviderEmail: "rafael.mendes.provider@petservices.local",
			ServiceName:   "Adestramento Básico",
			PetID:         "01KOWNP100000000000000003",
			Notes:         "Thor precisa melhorar obediência durante passeios.",
			Status:        "completed",
			RejectReason:  "",
		},
		{
			ID:            "01KOWNR100000000000000004",
			OwnerEmail:    "camila.nunes.owner@petservices.local",
			ProviderEmail: "mariana.costa.provider@petservices.local",
			ServiceName:   "Banho e Tosa Completo",
			PetID:         "01KOWNP100000000000000005",
			Notes:         "Bob tem pele sensível, usar shampoo hipoalergênico.",
			Status:        "completed",
			RejectReason:  "",
		},
		{
			ID:            "01KOWNR100000000000000005",
			OwnerEmail:    "camila.nunes.owner@petservices.local",
			ProviderEmail: "mariana.costa.provider@petservices.local",
			ServiceName:   "Corte de Unhas",
			PetID:         "01KOWNP100000000000000004",
			Notes:         "A Lua fica ansiosa, favor condução tranquila.",
			Status:        "rejected",
			RejectReason:  "Agenda indisponível para o horário solicitado.",
		},
		{
			ID:            "01KOWNR100000000000000006",
			OwnerEmail:    "diego.matos.owner@petservices.local",
			ProviderEmail: "carlos.eduardo.provider@petservices.local",
			ServiceName:   "Consulta Veterinária",
			PetID:         "01KOWNP100000000000000006",
			Notes:         "Consulta de rotina para check-up anual.",
			Status:        "completed",
			RejectReason:  "",
		},
		{
			ID:            "01KOWNR100000000000000007",
			OwnerEmail:    "diego.matos.owner@petservices.local",
			ProviderEmail: "carlos.eduardo.provider@petservices.local",
			ServiceName:   "Vacinação",
			PetID:         "01KOWNP100000000000000006",
			Notes:         "Atualização das vacinas obrigatórias.",
			Status:        "accepted",
			RejectReason:  "",
		},
		{
			ID:            "01KOWNR100000000000000008",
			OwnerEmail:    "paulo.dias.owner@petservices.local",
			ProviderEmail: "fernanda.lima.provider@petservices.local",
			ServiceName:   "Hospedagem (Hotel Pet)",
			PetID:         "01KOWNP100000000000000003",
			Notes:         "Reserva para fim de semana prolongado.",
			Status:        "pending",
			RejectReason:  "",
		},
	}

	providerServiceCache := map[string]models.Service{}

	affectedProviderIDs := map[string]struct{}{}
	for _, seed := range requestSeeds {
		owner, ok := ownersByEmail[seed.OwnerEmail]
		if !ok {
			return gorm.ErrRecordNotFound
		}

		provider, ok := providersByEmail[seed.ProviderEmail]
		if !ok {
			return gorm.ErrRecordNotFound
		}

		pet, ok := petsByID[seed.PetID]
		if !ok {
			return gorm.ErrRecordNotFound
		}

		serviceCacheKey := seed.ProviderEmail + "|" + seed.ServiceName
		service, found := providerServiceCache[serviceCacheKey]
		if !found {
			if err := db.Where("provider_id = ? AND name = ?", provider.ID, seed.ServiceName).First(&service).Error; err != nil {
				return err
			}
			providerServiceCache[serviceCacheKey] = service
		}

		request := models.Request{
			ID:           seed.ID,
			UserID:       owner.ID,
			ProviderID:   provider.ID,
			ServiceID:    service.ID,
			PetID:        pet.ID,
			Notes:        seed.Notes,
			Status:       seed.Status,
			RejectReason: seed.RejectReason,
			Active:       true,
		}

		if err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"user_id", "provider_id", "service_id", "pet_id", "notes", "status", "reject_reason", "active",
			}),
		}).Create(&request).Error; err != nil {
			return err
		}

		affectedProviderIDs[provider.ID] = struct{}{}
	}

	type reviewSeed struct {
		ID            string
		OwnerEmail    string
		ProviderEmail string
		Rating        float64
		Comment       string
	}

	reviewSeeds := []reviewSeed{
		{
			ID:            "01KOWNV100000000000000001",
			OwnerEmail:    "juliana.ferreira.owner@petservices.local",
			ProviderEmail: "ana.beatriz.provider@petservices.local",
			Rating:        4.8,
			Comment:       "Excelente cuidado com o Bento. Comunicação rápida e muito carinho com os pets.",
		},
		{
			ID:            "01KOWNV100000000000000002",
			OwnerEmail:    "paulo.dias.owner@petservices.local",
			ProviderEmail: "rafael.mendes.provider@petservices.local",
			Rating:        4.9,
			Comment:       "Treinamento muito eficiente. Thor evoluiu bastante em poucas sessões.",
		},
		{
			ID:            "01KOWNV100000000000000003",
			OwnerEmail:    "camila.nunes.owner@petservices.local",
			ProviderEmail: "mariana.costa.provider@petservices.local",
			Rating:        4.7,
			Comment:       "Ótimo atendimento, equipe cuidadosa e serviço de banho excelente.",
		},
		{
			ID:            "01KOWNV100000000000000004",
			OwnerEmail:    "diego.matos.owner@petservices.local",
			ProviderEmail: "carlos.eduardo.provider@petservices.local",
			Rating:        5.0,
			Comment:       "Consulta muito completa, explicações claras e ótimo acolhimento.",
		},
	}

	for _, seed := range reviewSeeds {
		owner, ok := ownersByEmail[seed.OwnerEmail]
		if !ok {
			return gorm.ErrRecordNotFound
		}

		provider, ok := providersByEmail[seed.ProviderEmail]
		if !ok {
			return gorm.ErrRecordNotFound
		}

		review := models.Review{
			ID:         seed.ID,
			UserID:     owner.ID,
			ProviderID: provider.ID,
			Rating:     seed.Rating,
			Comment:    seed.Comment,
			Active:     true,
		}

		if err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"user_id", "provider_id", "rating", "comment", "active",
			}),
		}).Create(&review).Error; err != nil {
			return err
		}

		affectedProviderIDs[provider.ID] = struct{}{}
	}

	for providerID := range affectedProviderIDs {
		var avgRating float64
		if err := db.Model(&models.Review{}).
			Where("provider_id = ? AND active = ?", providerID, true).
			Select("COALESCE(AVG(rating), 0)").
			Scan(&avgRating).Error; err != nil {
			return err
		}

		if err := db.Model(&models.Provider{}).
			Where("id = ?", providerID).
			Update("average_rating", avgRating).Error; err != nil {
			return err
		}
	}

	return nil
}

func Migration20260315000002(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Breed{}, &models.Pet{}); err != nil {
		return err
	}

	breedSeeds := []models.Breed{
		{ID: "01KBS000000000000000000001", Name: "Vira-lata", SpeciesID: "01KG7BG1XN0BQ4KHKPHY5V5ZEW", Active: true},
		{ID: "01KBS000000000000000000002", Name: "Labrador", SpeciesID: "01KG7BG1XN0BQ4KHKPHY5V5ZEW", Active: true},
		{ID: "01KBS000000000000000000003", Name: "Golden Retriever", SpeciesID: "01KG7BG1XN0BQ4KHKPHY5V5ZEW", Active: true},
		{ID: "01KBS000000000000000000004", Name: "Poodle", SpeciesID: "01KG7BG1XN0BQ4KHKPHY5V5ZEW", Active: true},
		{ID: "01KBS000000000000000000005", Name: "Shih Tzu", SpeciesID: "01KG7BG1XN0BQ4KHKPHY5V5ZEW", Active: true},
		{ID: "01KBS000000000000000000006", Name: "Bulldog Francês", SpeciesID: "01KG7BG1XN0BQ4KHKPHY5V5ZEW", Active: true},
		{ID: "01KBS000000000000000000007", Name: "Pastor Alemão", SpeciesID: "01KG7BG1XN0BQ4KHKPHY5V5ZEW", Active: true},
		{ID: "01KBS000000000000000000008", Name: "Beagle", SpeciesID: "01KG7BG1XN0BQ4KHKPHY5V5ZEW", Active: true},

		{ID: "01KBS000000000000000000009", Name: "Vira-lata", SpeciesID: "01KG7BG1XR2FA63J6N46CPVCKS", Active: true},
		{ID: "01KBS00000000000000000000A", Name: "Siamês", SpeciesID: "01KG7BG1XR2FA63J6N46CPVCKS", Active: true},
		{ID: "01KBS00000000000000000000B", Name: "Persa", SpeciesID: "01KG7BG1XR2FA63J6N46CPVCKS", Active: true},
		{ID: "01KBS00000000000000000000C", Name: "Maine Coon", SpeciesID: "01KG7BG1XR2FA63J6N46CPVCKS", Active: true},
		{ID: "01KBS00000000000000000000D", Name: "Angorá", SpeciesID: "01KG7BG1XR2FA63J6N46CPVCKS", Active: true},
		{ID: "01KBS00000000000000000000E", Name: "Ragdoll", SpeciesID: "01KG7BG1XR2FA63J6N46CPVCKS", Active: true},
		{ID: "01KBS00000000000000000000F", Name: "Bengal", SpeciesID: "01KG7BG1XR2FA63J6N46CPVCKS", Active: true},
		{ID: "01KBS00000000000000000000G", Name: "Sphynx", SpeciesID: "01KG7BG1XR2FA63J6N46CPVCKS", Active: true},
	}

	for _, breed := range breedSeeds {
		if err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "species_id"}, {Name: "name"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"active",
			}),
		}).Create(&breed).Error; err != nil {
			return err
		}
	}

	return nil
}
