use bevy::prelude::*;
use bevy_rapier2d::prelude::*;

use crate::components::{
    Aiming, BALL_MOVING_THRESHOLD, BALL_RADIUS, CUE_STICK_COLOR, CUE_STICK_GAP, CUE_STICK_LENGTH,
    CUE_STICK_PULLBACK, CUE_STICK_THICKNESS, CueBall, CueStick, MAX_IMPULSE, POWER_RATE, Power,
};
use crate::game::GamePhase;

pub fn spawn_cue_stick(
    mut commands: Commands,
    mut meshes: ResMut<Assets<Mesh>>,
    mut materials: ResMut<Assets<ColorMaterial>>,
) {
    commands.spawn((
        Mesh2d(meshes.add(Rectangle::new(CUE_STICK_LENGTH, CUE_STICK_THICKNESS))),
        MeshMaterial2d(materials.add(ColorMaterial::from(CUE_STICK_COLOR))),
        Transform::default(),
        Visibility::Visible,
        CueStick,
        Name::new("CueStick"),
    ));
}

#[allow(clippy::too_many_arguments, clippy::type_complexity)]
pub fn update_cue_aim(
    camera: Single<(&Camera, &GlobalTransform)>,
    window: Single<&Window>,
    mouse: Res<ButtonInput<MouseButton>>,
    time: Res<Time>,
    cue_ball: Single<(&Transform, &mut ExternalImpulse, &Velocity), With<CueBall>>,
    cue_stick: Single<(&mut Transform, &mut Visibility), (With<CueStick>, Without<CueBall>)>,
    mut aiming: ResMut<Aiming>,
    mut power: ResMut<Power>,
    mut next_phase: ResMut<NextState<GamePhase>>,
) {
    let (cam, cam_xf) = *camera;
    let Some(cursor) = window.cursor_position() else {
        return;
    };
    let Ok(cursor_world) = cam.viewport_to_world_2d(cam_xf, cursor) else {
        return;
    };
    aiming.cursor_world = cursor_world;

    let (ball_xf, mut impulse, vel) = cue_ball.into_inner();
    let aim = (cursor_world - ball_xf.translation.xy()).normalize_or_zero();
    aiming.aim = aim;
    let ball_moving = vel.linear.length_squared() > BALL_MOVING_THRESHOLD;

    let (mut stick_xf, mut stick_vis) = cue_stick.into_inner();

    if ball_moving {
        *stick_vis = Visibility::Hidden;
        aiming.charging = false;
        power.0 = 0.0;
        return;
    }
    *stick_vis = Visibility::Visible;

    if aim != Vec2::ZERO {
        let pull = CUE_STICK_GAP + power.0 * CUE_STICK_PULLBACK;
        let stick_center =
            ball_xf.translation.xy() - aim * (BALL_RADIUS + pull + CUE_STICK_LENGTH * 0.5);
        stick_xf.translation = stick_center.extend(0.7);
        stick_xf.rotation = Quat::from_rotation_z(aim.to_angle());
        stick_xf.scale = Vec3::ONE;
    }

    if mouse.just_pressed(MouseButton::Left) {
        aiming.charging = true;
        power.0 = 0.0;
    }
    if aiming.charging && mouse.pressed(MouseButton::Left) {
        power.0 = (power.0 + POWER_RATE * time.delta_secs()).min(1.0);
    }
    if mouse.just_released(MouseButton::Left) && aiming.charging {
        if aim != Vec2::ZERO {
            let shot = aim * (power.0 * MAX_IMPULSE);
            impulse.impulse = shot;

            // English: contact point offset = spin * R perpendicular to aim
            let perp = Vec2::new(-aim.y, aim.x);
            let offset = perp * aiming.spin * BALL_RADIUS;
            impulse.torque_impulse = offset.perp_dot(shot);

            next_phase.set(GamePhase::Shooting);
            *stick_vis = Visibility::Hidden;
        }
        power.0 = 0.0;
        aiming.charging = false;
    }
}
