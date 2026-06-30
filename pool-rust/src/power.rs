use bevy::prelude::*;

use crate::components::{
    Aiming, POWER_BAR_BG_COLOR, POWER_BAR_FILL_COLOR, POWER_BAR_HEIGHT, POWER_BAR_WIDTH,
    POWER_BAR_Y, Power, PowerBarBackground, PowerBarFill,
};

pub fn spawn_power_bar(
    mut commands: Commands,
    mut meshes: ResMut<Assets<Mesh>>,
    mut materials: ResMut<Assets<ColorMaterial>>,
) {
    commands.spawn((
        Mesh2d(meshes.add(Rectangle::new(POWER_BAR_WIDTH, POWER_BAR_HEIGHT))),
        MeshMaterial2d(materials.add(ColorMaterial::from(POWER_BAR_BG_COLOR))),
        Transform::from_xyz(0.0, POWER_BAR_Y, 0.8),
        Visibility::Hidden,
        PowerBarBackground,
        Name::new("PowerBarBackground"),
    ));

    commands.spawn((
        Mesh2d(meshes.add(Rectangle::new(1.0, POWER_BAR_HEIGHT))),
        MeshMaterial2d(materials.add(ColorMaterial::from(POWER_BAR_FILL_COLOR))),
        Transform::from_xyz(-POWER_BAR_WIDTH * 0.5, POWER_BAR_Y, 0.9),
        Visibility::Hidden,
        PowerBarFill,
        Name::new("PowerBarFill"),
    ));
}

pub fn update_power_bar(
    power: Res<Power>,
    aiming: Res<Aiming>,
    mut q_fill: Query<(&mut Transform, &mut Visibility), With<PowerBarFill>>,
    mut q_bg: Query<&mut Visibility, (With<PowerBarBackground>, Without<PowerBarFill>)>,
) {
    let show = aiming.charging;

    if let Ok(mut bg_vis) = q_bg.single_mut() {
        *bg_vis = if show {
            Visibility::Visible
        } else {
            Visibility::Hidden
        };
    }

    let Ok((mut fill_xf, mut fill_vis)) = q_fill.single_mut() else {
        return;
    };
    *fill_vis = if show {
        Visibility::Visible
    } else {
        Visibility::Hidden
    };
    if !show {
        return;
    }
    let w = (POWER_BAR_WIDTH * power.0).max(1.0);
    fill_xf.scale.x = w;
    fill_xf.scale.y = 1.0;
    fill_xf.translation.x = -POWER_BAR_WIDTH * 0.5 + w * 0.5;
    fill_xf.translation.y = POWER_BAR_Y;
}
