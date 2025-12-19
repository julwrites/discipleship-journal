import os
import cairosvg
from PIL import Image
import io

# Configuration
OUTPUT_DIR = "web/public"
ICON_COLOR = "#0f172a" # Slate 900
ACCENT_COLOR = "#ffffff" # White
PEN_COLOR = "#e2e8f0" # Slate 200 (lighter than accent)
BACKGROUND_COLOR = "#ffffff"

# SVG Templates
# A simple open book with a cross and a quill/pen
LOGO_SVG = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
  <defs>
    <linearGradient id="grad1" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" style="stop-color:#334155;stop-opacity:1" />
      <stop offset="100%" style="stop-color:#0f172a;stop-opacity:1" />
    </linearGradient>
  </defs>
  <!-- Background Circle for some contexts, but usually we want transparent or shaped.
       We will handle background in PNG generation if needed. -->

  <!-- Book Shape -->
  <path d="M464 96h-160c-26.51 0-48 21.49-48 48v288c0 26.51 21.49 48 48 48h160c26.51 0 48-21.49 48-48V144c0-26.51-21.49-48-48-48z" fill="{color}"/>
  <path d="M48 96h160c26.51 0 48 21.49 48 48v288c0 26.51-21.49 48-48 48H48c-26.51 0-48-21.49-48-48V144c0-26.51 21.49-48 48-48z" fill="{color}"/>

  <!-- Pages/Details (Accent) -->
  <path d="M416 160c0-8.84-7.16-16-16-16h-96c-8.84 0-16 7.16-16 16v224c0 8.84 7.16 16 16 16h96c8.84 0 16-7.16 16-16V160z" fill="{accent}" opacity="0.2"/>
  <path d="M96 160c0-8.84 7.16-16 16-16h96c8.84 0 16 7.16 16 16v224c0 8.84-7.16 16-16 16h-96c-8.84 0-16-7.16-16-16V160z" fill="{accent}" opacity="0.2"/>

  <!-- Cross in the middle -->
  <rect x="246" y="140" width="20" height="200" rx="4" fill="{accent}" />
  <rect x="180" y="180" width="152" height="20" rx="4" fill="{accent}" />

  <!-- Stylized Pen/Quill -->
  <!-- Diagonal placement on the right side -->
  <!-- A simple rectangle with a nib triangle -->
  <g transform="translate(320, 280) rotate(-45)">
    <!-- Pen Body -->
    <rect x="0" y="0" width="20" height="120" rx="2" fill="{pen_color}" />
    <!-- Pen Tip (Triangle) -->
    <path d="M0 120 L10 140 L20 120 Z" fill="{pen_color}" />
    <!-- Pen Cap/Detail -->
    <rect x="-2" y="-10" width="24" height="30" rx="2" fill="{accent}" />
  </g>

</svg>
"""

# Masked Icon (Monochrome black)
MASKED_SVG = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
  <path d="M464 96h-160c-26.51 0-48 21.49-48 48v288c0 26.51 21.49 48 48 48h160c26.51 0 48-21.49 48-48V144c0-26.51-21.49-48-48-48z" fill="black"/>
  <path d="M48 96h160c26.51 0 48 21.49 48 48v288c0 26.51-21.49 48-48 48H48c-26.51 0-48-21.49-48-48V144c0-26.51 21.49-48 48-48z" fill="black"/>
  <!-- Cross needs to be cutout (transparent) or handled differently for mask.
       For mask icon, usually it's a silhouette.
       Let's make the book solid and the cross transparent (cutout). -->
  <!-- Actually, for safari pinned tab, it uses the fill color. So we want the shape to be the ink.
       If we want the cross to be 'white' (background color), we should subtract it.
       Simplest is just the book shape if we can't do complex path ops easily here without a library.
       But we can overlay white on black in SVG. Mask icon treats opacity as nothing.
  -->
  <rect x="246" y="140" width="20" height="200" rx="4" fill="white" /> <!-- This won't work as 'cutout' in simple mask unless it's a mask. -->
  <!-- Let's just do a simple book shape without the cross cutout for the mask, or try to rely on path direction (nonzero winding) but that's complex to hand code.
       Alternative: Just the cross? No, book is better.
       Let's just use the book shape filled black. -->
</svg>
"""

def generate_svg(content, filename):
    with open(os.path.join(OUTPUT_DIR, filename), "w") as f:
        f.write(content)
    print(f"Generated {filename}")

def generate_png(svg_content, filename, width, height, bg_color=None):
    # If bg_color is provided, we might need to wrap the SVG or post-process with PIL
    # Cairosvg doesn't support 'background color' arg directly easily for transparent SVGs
    # unless we add a rect.

    svg_to_render = svg_content
    if bg_color:
        # Add a background rect
        bg_rect = f'<rect width="100%" height="100%" fill="{bg_color}"/>'
        # Insert after <svg ...>
        idx = svg_content.find(">") + 1
        svg_to_render = svg_content[:idx] + bg_rect + svg_content[idx:]

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

    # Prepare SVG content
    logo_svg = LOGO_SVG.format(color=ICON_COLOR, accent=ACCENT_COLOR, pen_color=PEN_COLOR)

    # 1. Save Base SVG
    generate_svg(logo_svg, "logo.svg")

    # 2. Generate PWA PNGs (Transparent background usually preferred for icons, but let's see.
    # Android icons often have transparent background.
    # iOS icons (apple-touch-icon) must be opaque (non-transparent).

    # PWA 192
    generate_png(logo_svg, "pwa-192x192.png", 192, 192)

    # PWA 512
    generate_png(logo_svg, "pwa-512x512.png", 512, 512)

    # OAuth Icon (High Res)
    generate_png(logo_svg, "oauth-icon.png", 1024, 1024)

    # Apple Touch Icon (Needs background)
    # We'll use a white background or the slate color with light icon?
    # Let's use white background with the slate icon.
    generate_png(logo_svg, "apple-touch-icon.png", 180, 180, bg_color="#ffffff")

    # Favicon (ICO) - base on 64x64 PNG
    # Use transparent background for favicon
    fav_png = generate_png(logo_svg, "favicon-source.png", 64, 64)
    generate_ico(fav_png, "favicon.ico")
    os.remove(os.path.join(OUTPUT_DIR, "favicon-source.png")) # Clean up temp

    # Masked Icon (SVG)
    # Use the monochrome version
    generate_svg(MASKED_SVG, "masked-icon.svg")

if __name__ == "__main__":
    main()
