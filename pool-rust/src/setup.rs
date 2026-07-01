use bevy::prelude::*;
use bevy_rapier2d::prelude::*;

use crate::components::*;

// (number, x, y) — bottom-to-top within each row (y ascending)
const RACK: &[(u8, f32, f32)] = &[
    // Row 0 (apex, foot spot)
    (1, 225.00, 0.00),
    // Row 1
    (2, 245.78, -12.00),
    (9, 245.78, 12.00),
    // Row 2 (8-ball dead center)
    (3, 266.57, -24.00),
    (8, 266.57, 0.00),
    (10, 266.57, 24.00),
    // Row 3
    (11, 287.35, -36.00),
    (4, 287.35, -12.00),
    (12, 287.35, 12.00),
    (13, 287.35, 36.00),
    // Row 4 (back two corners: 14 stripe / 7 solid — BCA rule)
    (14, 308.14, -48.00),
    (5, 308.14, -24.00),
    (15, 308.14, 0.00),
    (6, 308.14, 24.00),
    (7, 308.14, 48.00),
];

pub fn setup_physics(mut config: Query<&mut RapierConfiguration, With<DefaultRapierContext>>) {
    config.single_mut().unwrap().gravity = Vec2::ZERO;
}

pub fn spawn_camera(mut commands: Commands) {
    commands.spawn(Camera2d);
}

pub fn spawn_table(
    mut commands: Commands,
    mut meshes: ResMut<Assets<Mesh>>,
    mut materials: ResMut<Assets<ColorMaterial>>,
) {
    // Felt playing surface
    commands.spawn((
        Mesh2d(meshes.add(Rectangle::new(
            TABLE_HALF_WIDTH * 2.0,
            TABLE_HALF_HEIGHT * 2.0,
        ))),
        MeshMaterial2d(materials.add(ColorMaterial::from(FELT_COLOR))),
        Transform::from_xyz(0.0, 0.0, -1.0),
        Table,
    ));

    // Outer rail frame
    commands.spawn((
        Mesh2d(meshes.add(Rectangle::new(
            TABLE_HALF_WIDTH * 2.0 + CUSHION_THICKNESS * 2.0,
            TABLE_HALF_HEIGHT * 2.0 + CUSHION_THICKNESS * 2.0,
        ))),
        MeshMaterial2d(materials.add(ColorMaterial::from(RAIL_COLOR))),
        Transform::from_xyz(0.0, 0.0, -2.0),
    ));
}

pub fn spawn_cushions(mut commands: Commands) {
    let hw = TABLE_HALF_WIDTH;
    let hh = TABLE_HALF_HEIGHT;
    let cr = POCKET_RADIUS_CORNER;
    let sr = POCKET_RADIUS_SIDE;

    let segments: &[(Vec2, Vec2)] = &[
        // Top edge
        (Vec2::new(-hw + cr, hh), Vec2::new(-sr, hh)),
        (Vec2::new(sr, hh), Vec2::new(hw - cr, hh)),
        // Right edge
        (Vec2::new(hw, hh - cr), Vec2::new(hw, sr)),
        (Vec2::new(hw, -sr), Vec2::new(hw, -hh + cr)),
        // Bottom edge
        (Vec2::new(hw - cr, -hh), Vec2::new(sr, -hh)),
        (Vec2::new(-sr, -hh), Vec2::new(-hw + cr, -hh)),
        // Left edge
        (Vec2::new(-hw, -hh + cr), Vec2::new(-hw, -sr)),
        (Vec2::new(-hw, sr), Vec2::new(-hw, hh - cr)),
    ];

    for (a, b) in segments {
        commands.spawn((
            RigidBody::Fixed,
            Collider::segment(*a, *b),
            Transform::default(),
            Restitution {
                coefficient: CUSHION_RESTITUTION,
                combine_rule: CoefficientCombineRule::Min,
            },
            Friction::coefficient(CUSHION_FRICTION),
            Cushion,
        ));
    }
}

pub fn spawn_pockets(
    mut commands: Commands,
    mut meshes: ResMut<Assets<Mesh>>,
    mut materials: ResMut<Assets<ColorMaterial>>,
) {
    let pockets: &[(f32, f32, f32, PocketKind)] = &[
        (-450.0, 225.0, POCKET_RADIUS_CORNER, PocketKind::Corner),
        (0.0, 225.0, POCKET_RADIUS_SIDE, PocketKind::Side),
        (450.0, 225.0, POCKET_RADIUS_CORNER, PocketKind::Corner),
        (-450.0, -225.0, POCKET_RADIUS_CORNER, PocketKind::Corner),
        (0.0, -225.0, POCKET_RADIUS_SIDE, PocketKind::Side),
        (450.0, -225.0, POCKET_RADIUS_CORNER, PocketKind::Corner),
    ];

    for (x, y, radius, kind) in pockets {
        commands.spawn((
            Transform::from_xyz(*x, *y, 0.0),
            RigidBody::Fixed,
            Collider::ball(*radius),
            Sensor,
            ActiveCollisionTypes::all(),
            Pocket { kind: *kind },
            Mesh2d(meshes.add(Circle::new(*radius))),
            MeshMaterial2d(materials.add(ColorMaterial::from(POCKET_COLOR))),
        ));
    }
}

pub fn spawn_rack(
    mut commands: Commands,
    mut meshes: ResMut<Assets<Mesh>>,
    mut materials: ResMut<Assets<ColorMaterial>>,
) {
    for (n, x, y) in RACK {
        commands.spawn((
            Mesh2d(meshes.add(Circle::new(BALL_RADIUS))),
            MeshMaterial2d(materials.add(ColorMaterial::from(ball_color(*n)))),
            Transform::from_xyz(*x, *y, 0.0),
            RigidBody::Dynamic,
            Collider::ball(BALL_RADIUS),
            Restitution {
                coefficient: BALL_RESTITUTION,
                combine_rule: CoefficientCombineRule::Max,
            },
            Friction::coefficient(BALL_FRICTION),
            ActiveEvents::COLLISION_EVENTS,
            Damping {
                linear_damping: BALL_LINEAR_DAMPING,
                angular_damping: BALL_ANGULAR_DAMPING,
            },
            Sleeping::disabled(),
            Ball,
            Number { value: *n },
        ));
    }
}

pub fn spawn_cue_ball(
    mut commands: Commands,
    mut meshes: ResMut<Assets<Mesh>>,
    mut materials: ResMut<Assets<ColorMaterial>>,
) {
    commands.spawn((
        Mesh2d(meshes.add(Circle::new(BALL_RADIUS))),
        MeshMaterial2d(materials.add(ColorMaterial::from(ball_color(0)))),
        Transform::from_xyz(HEAD_SPOT_X, 0.0, 0.0),
        RigidBody::Dynamic,
        Collider::ball(BALL_RADIUS),
        Restitution {
            coefficient: BALL_RESTITUTION,
            combine_rule: CoefficientCombineRule::Max,
        },
        Friction::coefficient(BALL_FRICTION),
        ActiveEvents::COLLISION_EVENTS,
        Damping {
            linear_damping: BALL_LINEAR_DAMPING,
            angular_damping: BALL_ANGULAR_DAMPING,
        },
        Sleeping::disabled(),
        ExternalImpulse::default(),
        Ball,
        CueBall,
    ));
}
