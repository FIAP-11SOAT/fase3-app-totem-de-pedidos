model Category {
    id: string;
    name: string;
    description: string;
    created_at: string;
	updated_at : string;
}

model CategoryInput {
    name: string;
    description: string;
}