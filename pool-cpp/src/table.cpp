#include "table.h"

#include <cmath>

namespace {

void addSegment(std::vector<Segment>& segments, float x1, float y1, float x2, float y2) {
    segments.push_back({{x1, y1}, {x2, y2}});
}

void addCornerJaw(std::vector<Segment>& segments,
                  float cx,
                  float cy,
                  float jawLen,
                  float dx,
                  float dy) {
    addSegment(segments, cx, cy, cx + dx * jawLen, cy + dy * jawLen);
}

}  // namespace

void Table::build() {
    cushions.clear();
    pockets.clear();

    const float left = playArea.x;
    const float top = playArea.y;
    const float right = playArea.x + playArea.width;
    const float bottom = playArea.y + playArea.height;
    const float midX = playArea.x + playArea.width * 0.5f;

    const float pocketR = ballRadius * 1.55f;
    const float inset = pocketR * 0.82f;
    const float gap = pocketR * 1.28f;
    const float jawLen = pocketR * 0.72f;

    pockets = {
        {{left + inset, top + inset}, pocketR},
        {{midX, top - pocketR * 0.12f}, pocketR},
        {{right - inset, top + inset}, pocketR},
        {{left + inset, bottom - inset}, pocketR},
        {{midX, bottom + pocketR * 0.12f}, pocketR},
        {{right - inset, bottom - inset}, pocketR},
    };

    // Straight rails with pocket gaps.
    addSegment(cushions, left + gap, top, midX - gap, top);
    addSegment(cushions, midX + gap, top, right - gap, top);
    addSegment(cushions, left + gap, bottom, midX - gap, bottom);
    addSegment(cushions, midX + gap, bottom, right - gap, bottom);
    addSegment(cushions, left, top + gap, left, bottom - gap);
    addSegment(cushions, right, top + gap, right, bottom - gap);

    // Corner jaws guide balls into pockets instead of escaping through gaps.
    addCornerJaw(cushions, left + gap, top, jawLen, -1.0f, 1.0f);
    addCornerJaw(cushions, right - gap, top, jawLen, 1.0f, 1.0f);
    addCornerJaw(cushions, left + gap, bottom, jawLen, -1.0f, -1.0f);
    addCornerJaw(cushions, right - gap, bottom, jawLen, 1.0f, -1.0f);

    // Normalize jaw segment lengths to equal diagonal length.
    for (size_t i = cushions.size() - 4; i < cushions.size(); ++i) {
        Segment& seg = cushions[i];
        const Vec2 ab = seg.b - seg.a;
        const float len = ab.length();
        if (len > 1e-4f) {
            seg.b = seg.a + ab * (jawLen / len);
        }
    }
}
