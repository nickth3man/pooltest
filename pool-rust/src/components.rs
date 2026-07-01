#![allow(dead_code)]

use bevy::prelude::*;

pub const BALL_RADIUS: f32 = 12.0;
pub const TABLE_HALF_WIDTH: f32 = 450.0;
pub const TABLE_HALF_HEIGHT: f32 = 225.0;
pub const CUSHION_THICKNESS: f32 = 30.0;
pub const POCKET_GAP: f32 = 28.0;
pub const POCKET_RADIUS_CORNER: f32 = 16.0;
pub const POCKET_RADIUS_SIDE: f32 = 18.0;
pub const BALL_LINEAR_DAMPING: f32 = 0.02;
pub const BALL_ANGULAR_DAMPING: f32 = 0.2;
// Re-tuned restitution/friction
pub const BALL_RESTITUTION: f32 = 0.95;
pub const BALL_FRICTION: f32 = 0.05;
pub const CUSHION_RESTITUTION: f32 = 0.6;
pub const CUSHION_FRICTION: f32 = 0.3;
pub const PIXELS_PER_METER: f32 = 100.0;
pub const BALL_SPACING: f32 = 2.0 * BALL_RADIUS;
pub const ROW_SPACING: f32 = 20.784_609; // BALL_SPACING * (3.0_f32.sqrt() / 2.0) ≈ 20.7846
pub const HEAD_SPOT_X: f32 = -TABLE_HALF_WIDTH / 2.0;
pub const FOOT_SPOT_X: f32 = TABLE_HALF_WIDTH / 2.0;

pub const BG_COLOR: Color = Color::srgb(0.05, 0.35, 0.18);
pub const FELT_COLOR: Color = Color::srgb(0.06, 0.40, 0.20);
pub const RAIL_COLOR: Color = Color::srgb(0.18, 0.10, 0.05);
pub const POCKET_COLOR: Color = Color::BLACK;

pub const CUE_STICK_LENGTH: f32 = 240.0;
pub const CUE_STICK_THICKNESS: f32 = 6.0;
pub const CUE_STICK_GAP: f32 = 6.0;
pub const CUE_STICK_PULLBACK: f32 = 60.0;
pub const POWER_RATE: f32 = 1.2;
// Re-tuned shot power so max speed is ~5 m/s instead of ~11 m/s
pub const MAX_IMPULSE: f32 = 0.25;
pub const BALL_MOVING_THRESHOLD: f32 = 1.0;
pub const POWER_BAR_WIDTH: f32 = 220.0;
pub const POWER_BAR_HEIGHT: f32 = 16.0;
pub const POWER_BAR_Y: f32 = 350.0;
pub const CUE_STICK_COLOR: Color = Color::srgb(0.45, 0.27, 0.10);
pub const POWER_BAR_BG_COLOR: Color = Color::srgb(0.05, 0.05, 0.05);
pub const POWER_BAR_FILL_COLOR: Color = Color::srgb(0.95, 0.15, 0.10);
pub const CUT_LINE_LENGTH: f32 = 70.0;
pub const AIM_LINE_COLOR: Color = Color::srgba(1.0, 1.0, 1.0, 0.5);
pub const GHOST_BALL_COLOR: Color = Color::srgba(1.0, 1.0, 1.0, 0.45);

// Spin/English tuning
pub const SPIN_RATE: f32 = 1.0; // per second while holding ArrowLeft/ArrowRight or A/D

#[derive(Component)]
pub struct Ball;

#[derive(Component)]
pub struct CueBall;

#[derive(Component)]
pub struct CueStick;

#[derive(Component)]
pub struct PowerBarBackground;

#[derive(Component)]
pub struct PowerBarFill;

#[derive(Resource, Default)]
pub struct Power(pub f32);

#[derive(Resource, Default)]
pub struct Aiming {
    pub cursor_world: Vec2,
    pub charging: bool,
    pub aim: Vec2, // shared unit aim direction
    pub spin: f32, // -1..1 English, left/right
}

#[derive(Component)]
pub struct Number {
    pub value: u8,
}

#[derive(Component)]
pub struct Pocket {
    pub kind: PocketKind,
}

#[derive(Component)]
pub struct Cushion;

#[derive(Component)]
pub struct Table;

#[derive(Clone, Copy, Debug)]
pub enum PocketKind {
    Corner,
    Side,
}

#[derive(Clone, Copy, Debug)]
pub enum BallKind {
    Solid,
    Stripe,
    Eight,
}

pub fn ball_color(n: u8) -> Color {
    match n {
        1 | 9 => Color::srgb(1.000, 0.835, 0.000),  // yellow
        2 | 10 => Color::srgb(0.063, 0.165, 0.439), // dark blue
        3 | 11 => Color::srgb(0.847, 0.063, 0.094), // red
        4 | 12 => Color::srgb(0.357, 0.043, 0.467), // purple
        5 | 13 => Color::srgb(1.000, 0.498, 0.000), // orange
        6 | 14 => Color::srgb(0.063, 0.388, 0.122), // dark green
        7 | 15 => Color::srgb(0.439, 0.020, 0.122), // maroon
        8 => Color::srgb(0.043, 0.043, 0.043),      // black
        _ => Color::srgb(0.95, 0.95, 0.92),         // cue white (n=0)
    }
}

pub fn ball_kind(n: u8) -> BallKind {
    match n {
        8 => BallKind::Eight,
        1..=7 => BallKind::Solid,
        9..=15 => BallKind::Stripe,
        _ => BallKind::Solid, // cue defaults
    }
}
