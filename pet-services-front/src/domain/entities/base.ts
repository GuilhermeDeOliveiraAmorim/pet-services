export interface Base {
  id: string;
  active: boolean;
  createdAt: string;
  updatedAt?: string | null;
  deactivatedAt?: string | null;
}
