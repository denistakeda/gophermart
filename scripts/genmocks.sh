#!/usr/bin/env sh

mockgen -destination=mocks/core/ports/mockservice.go -package=ports gophermart/internal/core/ports \
    UserService,OrderService

mockgen -destination=mocks/core/ports/mockstore.go   -package=ports gophermart/internal/core/ports \
    UserStore,OrderStore,WithdrawnStore
