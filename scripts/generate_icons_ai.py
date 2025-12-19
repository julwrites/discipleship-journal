import os
import sys
from google import genai
from PIL import Image
import io

OUTPUT_DIR = "web/public"

def ensure_dir(directory):
    if not os.path.exists(directory):
        os.makedirs(directory)

def generate_icon():
    print("Initializing GenAI client...")
    try:
        # Client will automatically pick up GOOGLE_API_KEY from environment
        client = genai.Client()
        print("Client initialized.")

        prompt = (
            "A clean, minimal, modern mobile app icon for a discipleship journal app. "
            "A stylized open book with a simple cross. "
            "Deep slate blue background (#0f172a), with gradients of sky blue (#38bdf8) to indigo (#818cf8). "
            "Vector style, flat design, high fidelity, rounded corners, no text."
        )

        print(f"Generating image with prompt: {prompt}")
        response = client.models.generate_content(
            model="gemini-2.5-flash-image",
            contents=prompt,
        )

        if response.candidates and response.candidates[0].content.parts:
            for part in response.candidates[0].content.parts:
                if part.inline_data:
                    print("Image generated successfully.")
                    img_data = part.inline_data.data

                    # Save raw image
                    raw_path = os.path.join(OUTPUT_DIR, "logo-gen.png")
                    with open(raw_path, "wb") as f:
                        f.write(img_data)
                    print(f"Saved raw generated image to {raw_path}")
                    return raw_path
        else:
            print("No image content in response.")
            print(response)
            return None

    except Exception as e:
        print(f"Error generating icon: {e}")
        return None

def process_icons(source_path):
    if not source_path:
        print("No source image to process.")
        return

    print("Processing icons...")
    try:
        img = Image.open(source_path)

        # Ensure high quality resampling
        resample = Image.Resampling.LANCZOS

        # PWA Icons
        img.resize((192, 192), resample).save(os.path.join(OUTPUT_DIR, "pwa-192x192.png"))
        print("Generated pwa-192x192.png")

        img.resize((512, 512), resample).save(os.path.join(OUTPUT_DIR, "pwa-512x512.png"))
        print("Generated pwa-512x512.png")

        # Apple Touch Icon
        img.resize((180, 180), resample).save(os.path.join(OUTPUT_DIR, "apple-touch-icon.png"))
        print("Generated apple-touch-icon.png")

        # OAuth Icon (Used for consent screens, needs to be large)
        # If the generated image is smaller than 1024, resizing up might be blurry, but Gemini usually gives high res.
        # If it's already large or we are downscaling, it's fine.
        img.resize((1024, 1024), resample).save(os.path.join(OUTPUT_DIR, "oauth-icon.png"))
        print("Generated oauth-icon.png")

        # Favicon
        # .ico format can hold multiple sizes
        img.save(os.path.join(OUTPUT_DIR, "favicon.ico"), format='ICO', sizes=[(16, 16), (32, 32), (48, 48), (64, 64)])
        print("Generated favicon.ico")

    except Exception as e:
        print(f"Error processing icons: {e}")

if __name__ == "__main__":
    ensure_dir(OUTPUT_DIR)

    # Check if we want to skip generation and just process an existing file
    # useful for debugging or if we like the previous generation
    if len(sys.argv) > 1 and sys.argv[1] == "--process-only":
         raw_path = os.path.join(OUTPUT_DIR, "logo-gen.png")
         if os.path.exists(raw_path):
             process_icons(raw_path)
         else:
             print("No existing logo-gen.png found.")
    else:
        raw_path = generate_icon()
        if raw_path:
            process_icons(raw_path)
