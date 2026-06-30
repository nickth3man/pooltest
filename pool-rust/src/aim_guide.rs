use bevy::math::Isometry2d;
use bevy::prelude::*;
use bevy_rapier2d::prelude::*;

use crate::components::{
    AIM_LINE_COLOR, Aiming, BALL_MOVING_THRESHOLD, BALL_RADIUS, Ball, CUT_LINE_LENGTH, CueBall,
    GHOST_BALL_COLOR, Number, TABLE_HALF_HEIGHT, TABLE_HALF_WIDTH, ball_color,
};

/// Set the global gizmo line width once at startup.
pub fn set_aim_guide_width(mut config_store: ResMut<GizmoConfigStore>) {
    let (config, _) = config_store.config_mut::<DefaultGizmoConfigGroup>();
    config.line.width = 2.0;
}

#[allow(clippy::type_complexity)]
pub fn draw_aim_guide(
    mut gizmos: Gizmos,
    aiming: Res<Aiming>,
    cue_ball: Single<(&Transform, &Velocity), With<CueBall>>,
    balls: Query<(&Transform, &Number), (With<Ball>, Without<CueBall>)>,
) {
    // Only show when the cue ball is at rest (matches cue-stick visibility logic).
    let (ball_xf, vel) = *cue_ball;
    if vel.linear.length_squared() > BALL_MOVING_THRESHOLD {
        return;
    }

    let origin = ball_xf.translation.xy();
    let aim = (aiming.cursor_world - origin).normalize_or_zero();
    if aim == Vec2::ZERO {
        return;
    }

    // Find nearest object-ball contact along the ray.
    // Contact happens when cue-ball CENTER reaches distance 2*R from an object-ball center.
    let contact_radius = 2.0 * BALL_RADIUS;
    let mut best_t: Option<f32> = None;
    let mut hit_ball_pos: Option<Vec2> = None;
    let mut hit_ball_num: Option<u8> = None;
    for (obj_xf, num) in &balls {
        let p = obj_xf.translation.xy();
        let m = origin - p;
        let b = m.dot(aim);
        let c = m.dot(m) - contact_radius * contact_radius;
        let disc = b * b - c;
        if disc < 0.0 {
            continue;
        }
        let t = -b - disc.sqrt();
        if t <= 0.0 {
            continue;
        }
        if best_t.is_none_or(|best| t < best) {
            best_t = Some(t);
            hit_ball_pos = Some(p);
            hit_ball_num = Some(num.value);
        }
    }

    // Compute table-edge distance along the ray (slab method, AABB centered at world origin).
    let mut t_edge = f32::INFINITY;
    for (axis_origin, axis_dir, half) in [
        (origin.x, aim.x, TABLE_HALF_WIDTH),
        (origin.y, aim.y, TABLE_HALF_HEIGHT),
    ] {
        if axis_dir.abs() < 1e-6 {
            // Ray parallel to this slab; origin is inside the table so this wall is unreachable.
            continue;
        }
        let t1 = (-half - axis_origin) / axis_dir;
        let t2 = (half - axis_origin) / axis_dir;
        let texit = t1.max(t2);
        if texit < t_edge {
            t_edge = texit;
        }
    }

    // Pick the nearest obstruction: ball contact vs table edge.
    let t_ball = best_t.unwrap_or(f32::INFINITY);
    let t = t_ball.min(t_edge);
    if !t.is_finite() || t <= BALL_RADIUS {
        return;
    }

    // Draw the aim line from the cue-ball surface to the contact/edge point.
    let start = origin + aim * BALL_RADIUS;
    let end = origin + aim * t;
    gizmos.line_2d(start, end, AIM_LINE_COLOR);

    // If an object ball was struck, draw the ghost ball + the object ball's cut direction.
    if let (Some(t_hit), Some(obj_pos), Some(num)) = (best_t, hit_ball_pos, hit_ball_num) {
        if t_hit <= t_edge {
            let contact = origin + aim * t_hit;
            gizmos.circle_2d(
                Isometry2d::from_translation(contact),
                BALL_RADIUS,
                GHOST_BALL_COLOR,
            );
            // Object ball travels along the impact normal (contact -> object center), normalized.
            let n = (obj_pos - contact).normalize_or_zero();
            if n != Vec2::ZERO {
                gizmos.line_2d(obj_pos, obj_pos + n * CUT_LINE_LENGTH, ball_color(num));
            }
        }
    }
}
