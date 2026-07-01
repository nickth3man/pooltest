use bevy::prelude::*;
use bevy_rapier2d::parry::shape::SharedShape;
use bevy_rapier2d::prelude::*;

use crate::components::CueBall;
use crate::components::{
    Aiming, BALL_MOVING_THRESHOLD, BALL_RADIUS, Ball, Power, SPIN_RATE, TABLE_HALF_HEIGHT,
    TABLE_HALF_WIDTH,
};

#[derive(States, Debug, Clone, Copy, PartialEq, Eq, Hash, Default)]
pub enum GamePhase {
    #[default]
    Aiming,
    Shooting,
    BallInHand,
    #[allow(dead_code)]
    GameOver,
}

#[derive(Resource, Default)]
pub struct ShotResult {
    pub scratch: bool,
}

pub fn spin_input(
    keyboard: Res<ButtonInput<KeyCode>>,
    time: Res<Time>,
    mut aiming: ResMut<Aiming>,
) {
    if keyboard.pressed(KeyCode::ArrowLeft) || keyboard.pressed(KeyCode::KeyA) {
        aiming.spin = (aiming.spin - SPIN_RATE * time.delta_secs()).clamp(-1.0, 1.0);
    }
    if keyboard.pressed(KeyCode::ArrowRight) || keyboard.pressed(KeyCode::KeyD) {
        aiming.spin = (aiming.spin + SPIN_RATE * time.delta_secs()).clamp(-1.0, 1.0);
    }
}

pub fn check_balls_stopped(
    balls: Query<&Velocity, With<Ball>>,
    mut shot_result: ResMut<ShotResult>,
    mut power: ResMut<Power>,
    mut next_phase: ResMut<NextState<GamePhase>>,
) {
    let all_stopped = balls
        .iter()
        .all(|v| v.linear.length_squared() < BALL_MOVING_THRESHOLD);
    if all_stopped {
        if shot_result.scratch {
            next_phase.set(GamePhase::BallInHand);
        } else {
            next_phase.set(GamePhase::Aiming);
        }
        shot_result.scratch = false;
        power.0 = 0.0;
    }
}

#[allow(clippy::too_many_arguments, clippy::type_complexity)]
pub fn ball_in_hand(
    mut cue_ball: Query<(Entity, &mut Transform, &mut Velocity), With<CueBall>>,
    windows: Query<&Window>,
    camera: Query<(&Camera, &GlobalTransform)>,
    rapier_context: ReadRapierContext,
    mouse: Res<ButtonInput<MouseButton>>,
    mut next_phase: ResMut<NextState<GamePhase>>,
) {
    let Ok((cue_entity, mut transform, mut velocity)) = cue_ball.single_mut() else {
        return;
    };
    let Ok(window) = windows.single() else {
        return;
    };
    let Ok((cam, cam_gt)) = camera.single() else {
        return;
    };
    let Some(cursor) = window.cursor_position() else {
        return;
    };
    let Ok(world_pos) = cam.viewport_to_world_2d(cam_gt, cursor) else {
        return;
    };

    let half_w = TABLE_HALF_WIDTH - BALL_RADIUS;
    let half_h = TABLE_HALF_HEIGHT - BALL_RADIUS;
    let pos = world_pos.clamp(Vec2::new(-half_w, -half_h), Vec2::new(half_w, half_h));
    transform.translation = pos.extend(0.0);
    velocity.linear = Vec2::ZERO;
    velocity.angular = 0.0;

    let ctx = rapier_context.single().unwrap();
    let shape = SharedShape::ball(BALL_RADIUS);
    let mut overlaps = false;
    ctx.intersect_shape(
        pos,
        0.0,
        &*shape,
        QueryFilter::default()
            .exclude_collider(cue_entity)
            .exclude_sensors(),
        |_entity| {
            overlaps = true;
            false
        },
    );

    if mouse.just_pressed(MouseButton::Left) && !overlaps {
        next_phase.set(GamePhase::Aiming);
    }
}
