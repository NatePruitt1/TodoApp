## Plan: Project Domain Service/Handler Rescope

Keep one large ProjectService implementation for now (per your preference), but tighten it with clear sub-domain interfaces and complete all required project/category/card operations. Fix current compile blockers first, then implement missing methods and handlers with strict ownership checks so repository/service/handler responsibilities stay clean and predictable.

**Steps**
1. Phase 1: Stabilize compile baseline and interface consistency.
2. Fix naming/signature mismatches so wiring compiles: handler method name mismatch in main routing (GetAllProjectsHandler vs GetProjectsHandler), misspelled handler interface method (RemameCategory), invalid repository method usage in category service (GetCategory vs GetCategoryByID), and recursive DeleteCategory implementation.
3. Define complete request DTO set for project/category/card actions in request_dtos.go and align handler bindings to those DTOs. Keep DTO validation rules in handlers only.
4. Phase 2: Keep one large service but formalize boundaries internally.
5. Split ProjectService interface into grouped interface sections (ProjectOps, CategoryOps, CardOps) while retaining a single ProjectServiceImpl concrete type to avoid heavy wiring churn. This gives near-term simplicity with explicit responsibility boundaries.
6. Add required project operations for MVP completeness: update project name/description with per-user uniqueness policy consistent with AddProject.
7. Implement missing category operations in service: add, rename, delete, move/reorder with index normalization and bounds handling.
8. Implement missing card operations in service: add, rename, edit content, delete, move between categories, finish/unfinish toggle.
9. For every category/card mutating service method, enforce ownership via project owner checks using userId passed from handlers. Service owns authorization and business invariants; repository remains persistence-only.
10. Phase 3: Handler surface completion and route plan.
11. Implement project handler methods for update project and existing get/add/delete/get-one flows with consistent status codes and error mapping.
12. Implement category handler methods: add, rename, delete, move. Parse IDs from path params where possible and body for payload fields (name/index).
13. Implement card handler methods: add, rename, edit, delete, move, finish/unfinish.
14. Update cmd/api/main.go route registration to expose all required project/category/card endpoints under authenticated project group and remove stale handler names.
15. Phase 4: Complexity reduction and responsibility cleanup.
16. Replace insert-then-update patterns in project/category/card save methods with PostgreSQL upsert style for clarity and fewer round-trips.
17. Review aggregate fetching strategy: keep current GetAggregate for now, but isolate future optimization note (single joined query or batched query strategy) to avoid premature complexity.
18. Normalize error translation boundaries: repositories return technical errors only; services return domain errors (ownership, not found, conflict, invalid move); handlers map domain errors to HTTP codes.
19. Phase 5: Verification.
20. Add/expand service tests for ownership enforcement, duplicate-name checks, category reordering, card move across categories, and finish toggle behavior.
21. Add handler tests for DTO validation failures, auth context missing/invalid, and success payload contracts.
22. Run full backend build/tests and ensure no interface compile failures remain.

**Relevant files**
- /workspaces/TodoApp/apps/backend/internal/services/project_service.go — keep core service, group interface sections, add missing project/category/card methods.
- /workspaces/TodoApp/apps/backend/internal/services/category_service.go — merge/fold partial category logic into corrected ProjectServiceImpl methods.
- /workspaces/TodoApp/apps/backend/internal/handlers/project_handlers.go — implement missing handlers, fix naming typo, align with DTOs and status codes.
- /workspaces/TodoApp/apps/backend/internal/repository/project_repository.go — keep persistence-only API, simplify save methods, ensure method names used by services exist.
- /workspaces/TodoApp/apps/backend/internal/dto/request_dtos.go — define category/card/project update request DTOs.
- /workspaces/TodoApp/apps/backend/internal/dto/project_dtos.go — align response/request types and remove duplicate/unused patterns.
- /workspaces/TodoApp/apps/backend/internal/dto/mapper.go — resolve ProjectListItem mismatch and keep projection helpers coherent with DTOs.
- /workspaces/TodoApp/apps/backend/cmd/api/main.go — register final endpoint set and correct handler method names.

**Verification**
1. Build check: run backend compile to confirm all interfaces are satisfied and router method names exist.
2. Unit tests: execute service tests covering project update, category reorder edge cases, card move validations, and finish toggle.
3. Handler/API tests: validate auth failures, bad payload handling, and expected HTTP status codes for each new endpoint.
4. Regression check: verify existing create/login/refresh and existing project list/create/get/delete flows still pass.

**Decisions**
- Keep one large ProjectService implementation now, but structure by grouped interfaces and method sections.
- Include in MVP: category CRUD + reorder, card CRUD + move, card finish/unfinish, project update.
- Enforce ownership in service layer for all category/card mutations (handlers should not own authorization business rules beyond extracting authenticated user).
- Exclude from this iteration: deep query optimization rewrite of GetAggregate; only document and isolate for later.

**Further Considerations**
1. Endpoint style consistency: choose path-parameter-centric routes (recommended) versus request-body IDs for all operations; avoid mixing both patterns.
2. Reorder semantics: define whether move requests use absolute index only or support relative movement to reduce API ambiguity.
3. Domain error model: introducing typed sentinel errors in services will simplify clean HTTP mapping in handlers.