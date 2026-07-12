export interface UserRequestDTO {
    success: string;
    data: {
        id: string;
        username: string;
        created_at: Date;
        last_login?: Date;
    }
}