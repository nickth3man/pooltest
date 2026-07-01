#include "render.h"

#include "physics.h"

#include <cmath>
#include <cstdio>

namespace {

void drawBallBase(Image& image, int diameter, Color baseColor, bool stripe, int number) {
    const int cx = diameter / 2;
    const int cy = diameter / 2;
    const int radius = diameter / 2 - 1;

    ImageClearBackground(&image, BLANK);

    for (int y = 0; y < diameter; ++y) {
        for (int x = 0; x < diameter; ++x) {
            const float dx = static_cast<float>(x - cx);
            const float dy = static_cast<float>(y - cy);
            const float dist = std::sqrt(dx * dx + dy * dy);
            if (dist > radius) continue;

            Color pixel = baseColor;

            // Gloss highlight.
            const float highlight = std::exp(-((x - cx + radius * 0.25f) * (x - cx + radius * 0.25f) +
                                               (y - cy + radius * 0.35f) * (y - cy + radius * 0.35f)) /
                                              (radius * radius * 0.35f));
            pixel.r = static_cast<unsigned char>(std::min(255.0f, pixel.r + highlight * 70.0f));
            pixel.g = static_cast<unsigned char>(std::min(255.0f, pixel.g + highlight * 70.0f));
            pixel.b = static_cast<unsigned char>(std::min(255.0f, pixel.b + highlight * 70.0f));

            if (stripe && std::fabs(dy) < radius * 0.28f) {
                pixel = WHITE;
            }

            ImageDrawPixel(&image, x, y, pixel);
        }
    }

    if (number > 0) {
        char label[4];
        std::snprintf(label, sizeof(label), "%d", number);
        const int fontSize = diameter / 3;
        const int textW = MeasureText(label, fontSize);
        const Color textColor = (number == 8) ? WHITE : BLACK;
        ImageDrawText(&image, label, cx - textW / 2, cy - fontSize / 2, fontSize, textColor);
    }
}

}  // namespace

void BallSprites::create(int diameter) {
    destroy();
    size = diameter;

    for (int id = 0; id < kBallCount; ++id) {
        Image image = GenImageColor(diameter, diameter, BLANK);
        const bool stripe = id >= 9;
        const int number = (id == 0) ? 0 : id;
        drawBallBase(image, diameter, ballColor(id), stripe, number);
        textures[id] = LoadTextureFromImage(image);
        UnloadImage(image);
        SetTextureFilter(textures[id], TEXTURE_FILTER_BILINEAR);
    }
}

void BallSprites::destroy() {
    for (int id = 0; id < kBallCount; ++id) {
        if (textures[id].id != 0) {
            UnloadTexture(textures[id]);
            textures[id] = {};
        }
    }
    size = 0;
}

void BallSprites::draw(const Ball& ball) const {
    if (!ball.active || ball.id < 0 || ball.id >= kBallCount) return;

    const Texture2D& tex = textures[ball.id];
    const float d = static_cast<float>(size);
    const Rectangle src = {0.0f, 0.0f, d, d};
    const Rectangle dst = {
        ball.pos.x - ball.radius,
        ball.pos.y - ball.radius,
        ball.radius * 2.0f,
        ball.radius * 2.0f,
    };
    DrawTexturePro(tex, src, dst, {0.0f, 0.0f}, 0.0f, WHITE);
}

void Renderer::init() {
    sprites.create(static_cast<int>(Game::kBallRadius * 2.0f));
}

void Renderer::shutdown() {
    sprites.destroy();
}

void Renderer::draw(const Game& game) const {
    const Table& table = game.table;
    const Rectangle rail = {
        table.playArea.x - 18.0f,
        table.playArea.y - 18.0f,
        table.playArea.width + 36.0f,
        table.playArea.height + 36.0f,
    };

    DrawRectangleRec(rail, {55, 35, 20, 255});
    DrawRectangleRec(table.playArea, {34, 120, 60, 255});

    for (const Pocket& pocket : table.pockets) {
        DrawCircleV({pocket.center.x, pocket.center.y}, pocket.radius + 2.0f, BLACK);
        DrawCircleV({pocket.center.x, pocket.center.y}, pocket.radius, {15, 15, 15, 255});
    }

    for (const Segment& seg : table.cushions) {
        DrawLineEx({seg.a.x, seg.a.y}, {seg.b.x, seg.b.y}, 6.0f, {180, 140, 80, 255});
    }

    for (const Ball& ball : game.balls) {
        sprites.draw(ball);
    }

    const Ball* cue = game.cueBall();
    if (cue && game.canShoot() && game.aiming) {
        const Physics::RayHit hit = Physics::raycastFirstBall(
            cue->pos,
            game.aimDir,
            600.0f,
            game.balls,
            cue->radius,
            kCueBallId);

        const float lineLen =
            hit.hit ? hit.distance : (120.0f + game.aimPower * 80.0f);
        const Vec2 start = cue->pos;
        const Vec2 end = start + game.aimDir * lineLen;
        DrawLineEx({start.x, start.y}, {end.x, end.y}, 2.0f, Fade(RAYWHITE, 0.75f));

        if (hit.hit) {
            const Vec2 ghostCenter = hit.point;
            DrawCircleLines(ghostCenter.x, ghostCenter.y, cue->radius, Fade(SKYBLUE, 0.55f));
            DrawLineEx({start.x, start.y},
                       {ghostCenter.x, ghostCenter.y},
                       1.5f,
                       Fade(SKYBLUE, 0.35f));
        }

        const float cueLen = 140.0f + game.aimPower * 60.0f;
        const Vec2 cueEnd =
            start - game.aimDir * (cue->radius + 8.0f + game.aimPower * 40.0f);
        const Vec2 cueStart = cueEnd - game.aimDir * cueLen;
        DrawLineEx({cueStart.x, cueStart.y}, {cueEnd.x, cueEnd.y}, 5.0f, {210, 180, 120, 255});
        DrawCircleV({cueEnd.x, cueEnd.y}, 4.0f, {120, 80, 40, 255});
    }

    DrawText("Pool - C++", 16, 12, 24, RAYWHITE);
    DrawText("Aim with mouse | Pull back LMB to shoot | R to reset", 16, 40, 16, Fade(RAYWHITE, 0.85f));

    char status[64];
    if (!cue) {
        std::snprintf(status, sizeof(status), "Cue ball pocketed - press R to reset");
    } else if (game.ballsMoving) {
        std::snprintf(status, sizeof(status), "Balls in motion...");
    } else if (game.charging) {
        std::snprintf(status, sizeof(status), "Power: %.0f%%", game.aimPower * 100.0f);
    } else {
        std::snprintf(status, sizeof(status), "Balls on table: %d", game.activeBallCount());
    }
    DrawText(status, 16, 62, 16, Fade(RAYWHITE, 0.85f));

    DrawFPS(GetScreenWidth() - 100, 12);
}
