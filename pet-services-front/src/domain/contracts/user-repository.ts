import type { User } from "../entities";

export interface UserRepository {
  create(user: User): Promise<void>;
  findById(id: string): Promise<User | null>;
  findByEmail(email: string): Promise<User | null>;
  existsByEmail(email: string): Promise<boolean>;
  existsByPhone(countryCode: string, areaCode: string, number: string): Promise<boolean>;
  updateEmailVerified(userId: string, verified: boolean): Promise<void>;
  update(user: User): Promise<void>;
  delete(id: string): Promise<void>;
  list(page: number, limit: number): Promise<{ items: User[]; total: number }>;
}
