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

system "CollisionResolverSystem" {
    priority = 5
}

system "CollisionSystem" {
    priority = 6
}

system "GravitySystem" {
    priority = 7
}

entity "Camera" {
    component "Transform" {}
    component "Camera" {
        Zoom = 1.0
        # Bounds = {L: 0, B: 0, R: 800, T: 600}
    }
}

entity "TileMap" {
    component "Transform" {}
    component "Renderable" {}
    component "TileMap" {
        TileSize = 32
        Layer = 0
        Width = 25
        Height = 19
        Atlas = "assets/tiles/tileset.png"
        Tiles = []
    }
}