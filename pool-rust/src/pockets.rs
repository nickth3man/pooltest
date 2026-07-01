use bevy::prelude::*;
use bevy_rapier2d::prelude::*;
use bevy_rapier2d::rapier::geometry::CollisionEventFlags;

use crate::components::{Ball, CueBall, Pocket, TABLE_HALF_HEIGHT};
use crate::game::ShotResult;

pub fn detect_pocketed_balls(
    mut collision_events: MessageReader<CollisionEvent>,
    pockets: Query<&Pocket>,
    balls: Query<Entity, With<Ball>>,
    mut cue_balls: Query<(Entity, &mut Transform, &mut Velocity), With<CueBall>>,
    mut shot_result: ResMut<ShotResult>,
    mut commands: Commands,
) {
    for event in collision_events.read() {
        if let CollisionEvent::Started(a, b, flags) = event {
            if !flags.contains(CollisionEventFlags::SENSOR) {
                continue;
            }
            // resolve which entity is the pocket vs the ball
            let (_pocket_ent, ball_ent) = if pockets.get(*a).is_ok() && balls.get(*b).is_ok() {
                (*a, *b)
            } else if pockets.get(*b).is_ok() && balls.get(*a).is_ok() {
                (*b, *a)
            } else {
                continue;
            };
            // Scratch: move cue ball off-table to a holding position instead of despawning.
            if let Ok((_cue_entity, mut transform, mut velocity)) = cue_balls.get_mut(ball_ent) {
                transform.translation = Vec3::new(0.0, -TABLE_HALF_HEIGHT * 3.0, 0.0);
                velocity.linear = Vec2::ZERO;
                velocity.angular = 0.0;
                shot_result.scratch = true;
            } else {
                commands.entity(ball_ent).despawn();
            }
        }
    }
}
