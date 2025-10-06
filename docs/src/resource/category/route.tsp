import "@typespec/http";
import "@typespec/openapi";

import "./models.tsp";

using TypeSpec.Http;
using TypeSpec.OpenAPI;

namespace main.category;

@route("/categories")
@tag("Categories")
interface Category {
    
    @summary("Get all categories")
    @doc("Get all categories")
    @get ListAllCategories(): Category[];

    @summary("Get a categories by ID")
    @doc("Get a categories by ID")
    @get FindCategoryByID(@path id: string): Category;

    @summary("update a categories")
    @doc("update a categories")
    @put updateCategory(@path id: string, @body category: CategoryInput): void;

    @summary("Delete a categories")
    @doc("Delete a categories")
    @delete DeleteCategory(@path id: string): void;

    @summary("Create a new category")
    @doc("Create a new category")
    @post createCategory(@body category: CategoryInput): Category;

}