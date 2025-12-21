import os
import sys
import json

picdir = os.path.join(os.path.dirname(os.path.realpath(__file__)), 'RaspberryPi_JetsonNano/python/pic')
libdir = os.path.join(os.path.dirname(os.path.realpath(__file__)), 'RaspberryPi_JetsonNano/python/lib')
fontdir = os.path.join(os.path.dirname(os.path.realpath(__file__)), 'fonts')

if os.path.exists(libdir):
    sys.path.append(libdir)

from datetime import datetime
from PIL import Image, ImageDraw, ImageFont
from pathlib import Path

USE_FAKE_EPD = True

if not USE_FAKE_EPD:
    from waveshare_epd import epd2in7 as epd2in7

JSON_PATH = Path(__file__).parent.parent / "next.json"

class FakeEPD:
    width = 176
    height = 264

    def init(self):
        print("FakeEPD init")

    def Clear(self, color=255):
        print("FakeEPD clear")

    def display(self, image_buffer):
        # Convert buffer back to Image if needed
        print("FakeEPD display called")
        # image_buffer is already a PIL Image in Waveshare library
        image_buffer.show(title="Fake EPD Display")

    def getbuffer(self, image):
        # Just return the PIL Image itself
        return image

    def sleep(self):
        print("FakeEPD sleep")

def load_json():
    if not JSON_PATH.exists():
        raise FileNotFoundError(f"{JSON_PATH} does not exist")

    with JSON_PATH.open("r") as f:
        return json.load(f)


def parse_predictions(data):
    prds = data["Bustime"]["bustime-response"].get("prd", [])
    results = []

    for p in prds:
        results.append({
            "dist": p["dstp"],
            "route": p["rtdd"],
            "dir": p["rtdir"],
            "dest": p["des"],
            "mins": p["prdctdn"],
            "delayed": p["dly"],
        })

    return results


def draw_display(predictions, generated_at):
    epd = FakeEPD() if USE_FAKE_EPD else epd2in7.EPD()

    epd.init()
    epd.Clear(0xFF)

    image = Image.new("1", (epd.width, epd.height), 255)
    draw = ImageDraw.Draw(image)

    font_title = ImageFont.truetype(fontdir + "/DejaVuSansMono-Bold.ttf", 22)
    font_body = ImageFont.truetype(fontdir + "/DejaVuSansMono-Bold.ttf", 18)

    y = 0

    # Header
    draw.text((0, y), "Paterson", font=font_title, fill=0)
    y += 30

    # Predictions
    for p in predictions[:5]:
        dist = int(p["dist"])
        if dist == 0:
            distance = ""
        elif dist < 5280:
            distance = f"<1 mi"
        else:
            distance = f"{dist / 5280:.1f} mi"

        route_letter = p["route"]
        route_info = f"{p['mins']} min{distance}"

        # Draw black circle
        radius = 10
        cx = radius
        cy = y + radius
        draw.ellipse((cx - radius, cy - radius, cx + radius, cy + radius), fill=0, outline=0)

        # Draw route letter in the center as white text
        # anchor=mm means 'middle' horizontal, 'middle' vertical
        draw.text((cx, cy), route_letter, font=font_body, fill=255, anchor="mm")

        tx = 2 * radius + 4
        draw.text((tx, y), f"{p['mins']:>2} min", font=font_body, fill=0)
        if distance:
            dx = epd.width - 4
            # anchor=ra means 'right' horizontal, 'baseline' vertical
            draw.text((dx, y), distance, font=font_body, fill=0, anchor="ra")
        
        text_bbox = draw.textbbox((0, 0), route_info, font=font_body)
        font_height = text_bbox[3] - text_bbox[1]  # bottom - top
        row_height = max(2 * radius, font_height) + 4  # extra padding
        
        y += row_height

    y += 4
    draw.line((0, y, epd.width, y), fill=0)
    y += 6

    draw.text((0, y), f"Updated At {generated_at.strftime('%H:%M:%S')}", font=font_body, fill=0)

    epd.display(epd.getbuffer(image))
    epd.sleep()


def main():
    data = load_json()

    predictions = parse_predictions(data)
    generated_at = datetime.fromisoformat(data["GeneratedAt"])

    draw_display(predictions, generated_at)


if __name__ == "__main__":
    main()
