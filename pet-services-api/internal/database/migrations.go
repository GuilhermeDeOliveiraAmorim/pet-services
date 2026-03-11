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
		AverageRating: 4.8,
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
