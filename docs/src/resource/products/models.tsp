model Product {
    id: int32;
    name: string;
    description: string;
    price: float64;
    imageUrl: string;
    preparationTime: int32;
    createdAt: string;
    updatedAt: string;
    categoryId: int32;
}

model ProductInput {
    name: string;
    description: string;
    price: float64;
    imageUrl: string;
    categoryId: int32;
}