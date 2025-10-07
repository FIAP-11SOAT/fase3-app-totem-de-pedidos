model Customer {
    ID: string;
    Name: string;
    Email: string;
    TaxID: string;
    CreatedAt: string;
}

model CustomerInput {
    taxid: string;
    nome: string;
    email: string;
}