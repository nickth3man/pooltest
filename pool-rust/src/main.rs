//! Pool - Rust: Bevy 2D pool/billiards game.
// Support configuring Bevy lints within code when `bevy_lint` runs.
#![cfg_attr(bevy_lint, feature(register_tool), register_tool(bevy))]

use bevy::prelude::*;
use bevy_rapier2d::prelude::*;

#[cfg(feature = "debug")]
use bevy_rapier2d::render::RapierDebugRenderPlugin;

mod aim_guide;
mod components;
mod cue;
mod pockets;
mod power;
mod setup;

fn main() {
    let mut app = App::new();
    app.insert_resource(ClearColor(components::BG_COLOR))
        .add_plugins(DefaultPlugins.set(WindowPlugin {
            primary_window: Some(Window {
                resolution: (1024, 768).into(),
                title: "Pool".into(),
                ..default()
            }),
            ..default()
        }))
        .add_plugins(RapierPhysicsPlugin::<NoUserData>::pixels_per_meter(
            components::PIXELS_PER_METER,
        ))
        .insert_resource(components::Power::default())
        .insert_resource(components::Aiming::default())
        .add_systems(
            Startup,
            (
                setup::setup_physics,
                setup::spawn_camera,
                setup::spawn_table,
                setup::spawn_cushions,
                setup::spawn_pockets,
                setup::spawn_rack,
                setup::spawn_cue_ball,
                cue::spawn_cue_stick,
                power::spawn_power_bar,
                aim_guide::set_aim_guide_width,
            ),
        )
        .add_systems(
            Update,
            (
                cue::update_cue_aim,
                power::update_power_bar,
                aim_guide::draw_aim_guide,
            ),
        )
        .add_systems(PostUpdate, pockets::detect_pocketed_balls);

    #[cfg(feature = "debug")]
    app.add_plugins(RapierDebugRenderPlugin::default());

    app.run();
}
