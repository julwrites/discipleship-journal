import os
import cairosvg
from PIL import Image
import io

# Configuration
OUTPUT_DIR = "web/public"
# Using a "Nano Banana" inspired modern palette
# Deep Slate/Blue background, Vivid Gradient
ICON_BG_COLOR = "#0f172a" # Slate 900
ICON_ACCENT_GRADIENT_START = "#38bdf8" # Sky 400
ICON_ACCENT_GRADIENT_END = "#818cf8" # Indigo 400
WHITE = "#ffffff"

# SVG Templates

# Modern Icon:
# Rounded Square Background (handled by shape or container in some contexts, but here we bake it in for the icon)
# We will make the "Logo" just the symbol (Book + Cross), and put it on a background for the square icons.

# Symbol: A stylized open book with a cross in negative space or accent color.
# Let's do a book that looks like a shield/heart hybrid? No, stick to book.
# Modern, thick lines.

LOGO_SYMBOL = """
<g transform="translate(106, 106) scale(0.6)"> <!-- Center and scale -->
  <!-- Book Base -->
  <path d="M448 360V48a48 48 0 0 0-48-48H128a48 48 0 0 0-48 48v312c0 26.5 21.5 48 48 48h272a48 48 0 0 0 48-48z" fill="url(#gradSymbol)" />
  <path d="M128 0a48 48 0 0 0-48 48v312c0 26.5 21.5 48 48 48h272v-48H128a16 16 0 0 1-16-16V16a16 16 0 0 1 16-16h272V0H128z" fill="#000" opacity="0.1"/>

  <!-- Pages / Detail -->
  <path d="M400 48H128a16 16 0 0 0-16 16v320a16 16 0 0 0 16 16h272" fill="none" stroke="#fff" stroke-width="24" stroke-linecap="round"/>

  <!-- Cross Symbol in Center (White) -->
  <path d="M240 100v160M180 160h120" stroke="#fff" stroke-width="32" stroke-linecap="round" />

  <!-- Quill/Pen Accent (Gold/Yellow for pop) -->
  <path d="M380 320 l-60 60 l-20 -10 l10 -20 l60 -60 z" fill="#fbbf24" />
</g>
"""

# Full Icon with Background (for PWA, etc)
LOGO_SVG = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
  <defs>
    <linearGradient id="gradBg" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" style="stop-color:#1e293b;stop-opacity:1" />
      <stop offset="100%" style="stop-color:#0f172a;stop-opacity:1" />
    </linearGradient>
    <linearGradient id="gradSymbol" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" style="stop-color:#38bdf8;stop-opacity:1" />
      <stop offset="100%" style="stop-color:#818cf8;stop-opacity:1" />
    </linearGradient>
  </defs>

  <!-- Background Squircle -->
  <rect x="0" y="0" width="512" height="512" rx="128" fill="url(#gradBg)"/>

  <!-- Symbol -->
  """ + LOGO_SYMBOL + """
</svg>
"""

# Masked Icon (Monochrome Black)
MASKED_SVG = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
  <!-- Just the book shape in black -->
  <g transform="translate(106, 106) scale(0.6)">
    <path d="M448 360V48a48 48 0 0 0-48-48H128a48 48 0 0 0-48 48v312c0 26.5 21.5 48 48 48h272a48 48 0 0 0 48-48z" fill="black" />
  </g>
</svg>
"""

# Apple Touch Icon needs a specific background usually, but our LOGO_SVG has a background.
# We will use LOGO_SVG directly.

def generate_svg(content, filename):
    with open(os.path.join(OUTPUT_DIR, filename), "w") as f:
        f.write(content)
    print(f"Generated {filename}")

def generate_png(svg_content, filename, width, height, bg_color=None):
    svg_to_render = svg_content
    # If we need to force a background that isn't in the SVG (though our SVG has one now)
    # logic would go here. But we designed LOGO_SVG to be the full icon.

    png_data = cairosvg.svg2png(bytestring=svg_to_render.encode('utf-8'), output_width=width, output_height=height)

    with open(os.path.join(OUTPUT_DIR, filename), "wb") as f:
        f.write(png_data)
    print(f"Generated {filename}")
    return png_data

def generate_ico(png_data, filename):
    img = Image.open(io.BytesIO(png_data))
    img.save(os.path.join(OUTPUT_DIR, filename), format='ICO', sizes=[(16,16), (32,32), (48,48), (64,64)])
    print(f"Generated {filename}")

def main():
    if not os.path.exists(OUTPUT_DIR):
        os.makedirs(OUTPUT_DIR)

    # 1. Save Base SVG (The full icon)
    generate_svg(LOGO_SVG, "logo.svg")

    # 2. PWA Icons
    generate_png(LOGO_SVG, "pwa-192x192.png", 192, 192)
    generate_png(LOGO_SVG, "pwa-512x512.png", 512, 512)

    # 3. OAuth Icon
    generate_png(LOGO_SVG, "oauth-icon.png", 1024, 1024)

    # 4. Apple Touch Icon
    generate_png(LOGO_SVG, "apple-touch-icon.png", 180, 180)

    # 5. Favicon
    # For favicon, we might want just the symbol without the big background box?
    # Or just the box. Let's use the box for consistency.
    fav_png = generate_png(LOGO_SVG, "favicon-source.png", 64, 64)
    generate_ico(fav_png, "favicon.ico")
    os.remove(os.path.join(OUTPUT_DIR, "favicon-source.png"))

    # 6. Masked Icon
    generate_svg(MASKED_SVG, "masked-icon.svg")

if __name__ == "__main__":
    main()
