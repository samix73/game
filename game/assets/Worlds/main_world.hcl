name = "main_world"

system "CameraSystem" {
    priority = 1
}

system "PhysicsSystem" {
    priority = 2
}

system "PauseSystem" {
    priority = 3
}

system "TileSystem" {
    priority = 4
}

system "CollisionSystem" {
    priority = 5
}

system "GravitySystem" {
    priority = 6
}

