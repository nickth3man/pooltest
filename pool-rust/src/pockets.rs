use bevy::prelude::*;
use bevy_rapier2d::prelude::*;
use bevy_rapier2d::rapier::geometry::CollisionEventFlags;

use crate::components::{Ball, CueBall, HEAD_SPOT_X, Pocket};

pub fn detect_pocketed_balls(
    mut collision_events: MessageReader<CollisionEvent>,
    pockets: Query<&Pocket>,
    balls: Query<Entity, With<Ball>>,
    mut cue_balls: Query<(&mut Transform, &mut Velocity), With<CueBall>>,
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
            // Scratch: respawn cue ball at the head spot instead of despawning.
            if let Ok((mut transform, mut velocity)) = cue_balls.get_mut(ball_ent) {
                transform.translation = Vec3::new(HEAD_SPOT_X, 0.0, 0.0);
                velocity.linear = Vec2::ZERO;
                velocity.angular = 0.0;
            } else {
                commands.entity(ball_ent).despawn();
            }
        }
    }
}
