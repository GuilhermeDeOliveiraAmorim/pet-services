export interface LogoutInput {
  userId: string;
  tokenId?: string;
  revokeAll: boolean;
}

export interface LogoutInputBody {
  tokenId?: string;
  revokeAll: boolean;
}

export interface LogoutOutput {
  message?: string;
  detail?: string;
}
