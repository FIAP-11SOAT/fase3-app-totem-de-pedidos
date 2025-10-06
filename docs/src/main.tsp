import "@typespec/http";
import "@typespec/openapi";

import "./constants.tsp";
import "./resource/customers/route.tsp";
import "./resource/orders/route.tsp";
import "./resource/category/route.tsp";
import "./resource/products/route.tsp";
import "./resource/payment/route.tsp";

using TypeSpec.OpenAPI;
using TypeSpec.Http;


@service()
@info(Info)
@doc(DescMain)
namespace main;
