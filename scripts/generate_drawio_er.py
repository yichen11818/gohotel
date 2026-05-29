from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path
import shutil
import xml.etree.ElementTree as ET

from PIL import Image, ImageDraw, ImageFont


ROOT = Path(__file__).resolve().parents[1]
ASSET_DIR = ROOT / "output" / "doc" / "assets"
DRAWIO_OUT = ASSET_DIR / "er-redesign-drawio.drawio"
PNG_OUT = ASSET_DIR / "er-redesign-drawio.png"
RADIAL_DRAWIO_OUT = ASSET_DIR / "er-radial-drawio.drawio"
RADIAL_PNG_OUT = ASSET_DIR / "er-radial-drawio.png"

TITLE = "系统核心业务 E-R 图"
PAGE_W = 2200
PAGE_H = 1600


def load_font(size: int, bold: bool = False) -> ImageFont.FreeTypeFont:
    paths = [
        r"C:\Windows\Fonts\simhei.ttf" if bold else r"C:\Windows\Fonts\simsun.ttc",
        r"C:\Windows\Fonts\msyhbd.ttc" if bold else r"C:\Windows\Fonts\msyh.ttc",
    ]
    for item in paths:
        path = Path(item)
        if path.exists():
            return ImageFont.truetype(str(path), size)
    return ImageFont.load_default()


FONT_TITLE = load_font(38, True)
FONT_ENTITY = load_font(30, True)
FONT_ATTR = load_font(20)
FONT_REL = load_font(21, True)
FONT_CARD = load_font(18, True)
FONT_NOTE = load_font(18)


@dataclass(frozen=True)
class Node:
    id: str
    value: str
    kind: str
    x: int
    y: int
    w: int
    h: int

    @property
    def cx(self) -> int:
        return self.x + self.w // 2

    @property
    def cy(self) -> int:
        return self.y + self.h // 2


@dataclass(frozen=True)
class Attr:
    entity: str
    attr: str


@dataclass(frozen=True)
class Rel:
    a: str
    rel: str
    b: str
    a_card: str
    b_card: str
    a_label: tuple[int, int]
    b_label: tuple[int, int]


NODES: list[Node] = [
    Node("room", "房间", "entity", 1000, 755, 200, 90),
    Node("room_no", "房间号", "attr", 875, 600, 120, 46),
    Node("room_floor", "楼层", "attr", 1050, 560, 100, 46),
    Node("room_status", "房态", "attr", 1235, 610, 100, 46),
    Node("room_clean", "清洁状态", "attr", 1250, 870, 135, 46),
    Node("room_profile", "价格/容量", "attr", 850, 900, 145, 46),
    Node("room_pos", "位置尺寸", "attr", 1065, 940, 135, 46),

    Node("booking", "预订", "entity", 390, 735, 170, 66),
    Node("booking_no", "订单号", "attr", 290, 610, 120, 46),
    Node("booking_date", "入住/离店", "attr", 465, 585, 145, 46),
    Node("booking_guest", "入住人信息", "attr", 280, 875, 150, 46),
    Node("booking_status", "订单/支付状态", "attr", 510, 905, 170, 46),

    Node("user", "用户", "entity", 285, 300, 150, 66),
    Node("user_id", "用户ID", "attr", 95, 170, 115, 46),
    Node("user_account", "用户名", "attr", 290, 130, 120, 46),
    Node("user_role", "角色", "attr", 505, 185, 105, 46),
    Node("user_member", "等级/积分", "attr", 500, 380, 135, 46),

    Node("behavior", "用户行为", "entity", 290, 1065, 170, 66),
    Node("beh_type", "行为类型", "attr", 95, 1215, 130, 46),
    Node("beh_weight", "权重/来源", "attr", 390, 1230, 140, 46),

    Node("hotel", "酒店配置", "entity", 765, 250, 170, 66),
    Node("hotel_profile", "基础资料", "attr", 625, 110, 130, 46),
    Node("hotel_rule", "预订规则", "attr", 885, 110, 130, 46),

    Node("category", "房型分类", "entity", 1530, 300, 180, 66),
    Node("cat_name", "分类名", "attr", 1375, 155, 120, 46),
    Node("cat_facility", "设施", "attr", 1580, 110, 100, 46),
    Node("cat_image", "图片", "attr", 1790, 190, 100, 46),

    Node("operation", "运营内容", "entity", 1830, 495, 170, 66),
    Node("op_type", "公告/轮播", "attr", 1940, 345, 140, 46),
    Node("op_valid", "排序/有效期", "attr", 1990, 610, 150, 46),

    Node("inventory", "房态库存", "entity", 1765, 820, 180, 66),
    Node("inv_date", "日期", "attr", 1710, 650, 100, 46),
    Node("inv_count", "总量/已售量", "attr", 1940, 695, 160, 46),
    Node("inv_price", "当日价格", "attr", 1815, 1015, 135, 46),

    Node("pricing", "定价规则", "entity", 1570, 1190, 180, 66),
    Node("price_type", "规则类型", "attr", 1485, 1330, 130, 46),
    Node("price_adjust", "调整值/有效期", "attr", 1595, 1435, 165, 46),
    Node("price_priority", "优先级", "attr", 1835, 1335, 115, 46),

    Node("task", "房务任务", "entity", 850, 1210, 170, 66),
    Node("task_type", "清洁/维修", "attr", 665, 1335, 135, 46),
    Node("task_status", "状态", "attr", 885, 1440, 100, 46),
    Node("task_staff", "处理人", "attr", 1070, 1345, 115, 46),

    Node("facility", "楼层设施", "entity", 1230, 1215, 170, 66),
    Node("fac_type", "设施类型", "attr", 1195, 1445, 130, 46),
    Node("fac_pos", "位置尺寸", "attr", 1425, 1405, 130, 46),
]

REL_NODES: list[Node] = [
    Node("rel_room", "对应", "rel", 690, 765, 92, 74),
    Node("rel_submit", "提交", "rel", 410, 525, 92, 74),
    Node("rel_behavior_user", "产生", "rel", 300, 680, 92, 74),
    Node("rel_behavior_room", "指向", "rel", 620, 1000, 92, 74),
    Node("rel_hotel_booking", "约束", "rel", 650, 515, 92, 74),
    Node("rel_belong", "归属", "rel", 1320, 510, 92, 74),
    Node("rel_operation", "展示", "rel", 1725, 400, 92, 74),
    Node("rel_stock", "关联", "rel", 1565, 780, 92, 74),
    Node("rel_price", "作用于", "rel", 1420, 1040, 100, 78),
    Node("rel_task", "处理", "rel", 960, 1030, 92, 74),
    Node("rel_facility", "布局", "rel", 1210, 1035, 92, 74),
]

ATTRS: list[Attr] = [
    Attr("booking", "booking_no"),
    Attr("booking", "booking_date"),
    Attr("booking", "booking_guest"),
    Attr("booking", "booking_status"),
    Attr("user", "user_id"),
    Attr("user", "user_account"),
    Attr("user", "user_role"),
    Attr("user", "user_member"),
    Attr("room", "room_no"),
    Attr("room", "room_floor"),
    Attr("room", "room_status"),
    Attr("room", "room_clean"),
    Attr("room", "room_profile"),
    Attr("room", "room_pos"),
    Attr("category", "cat_name"),
    Attr("category", "cat_facility"),
    Attr("category", "cat_image"),
    Attr("inventory", "inv_date"),
    Attr("inventory", "inv_count"),
    Attr("inventory", "inv_price"),
    Attr("pricing", "price_type"),
    Attr("pricing", "price_adjust"),
    Attr("pricing", "price_priority"),
    Attr("task", "task_type"),
    Attr("task", "task_status"),
    Attr("task", "task_staff"),
    Attr("behavior", "beh_type"),
    Attr("behavior", "beh_weight"),
    Attr("hotel", "hotel_profile"),
    Attr("hotel", "hotel_rule"),
    Attr("operation", "op_type"),
    Attr("operation", "op_valid"),
    Attr("facility", "fac_type"),
    Attr("facility", "fac_pos"),
]

RELS: list[Rel] = [
    Rel("booking", "rel_room", "room", "N", "1", (590, 735), (950, 735)),
    Rel("user", "rel_submit", "booking", "1", "N", (450, 445), (460, 690)),
    Rel("user", "rel_behavior_user", "behavior", "1", "N", (320, 410), (330, 1035)),
    Rel("behavior", "rel_behavior_room", "room", "N", "1", (495, 1050), (960, 855)),
    Rel("hotel", "rel_hotel_booking", "booking", "1", "N", (775, 385), (570, 690)),
    Rel("room", "rel_belong", "category", "N", "1", (1220, 685), (1510, 410)),
    Rel("category", "rel_operation", "operation", "1", "N", (1715, 385), (1835, 455)),
    Rel("room", "rel_stock", "inventory", "1", "N", (1210, 780), (1730, 800)),
    Rel("room", "rel_price", "pricing", "1", "N", (1200, 865), (1560, 1145)),
    Rel("room", "rel_task", "task", "1", "N", (1060, 860), (955, 1175)),
    Rel("room", "rel_facility", "facility", "N", "N", (1200, 865), (1260, 1180)),
]

NOTE = "说明：房务任务包含清洁任务与维修工单；运营内容包含公告与轮播图；酒店配置包含基础资料与预订规则。"


def edge_point(a: Node, b: Node) -> tuple[int, int]:
    dx = b.cx - a.cx
    dy = b.cy - a.cy
    if abs(dx) / max(a.w, 1) > abs(dy) / max(a.h, 1):
        x = a.x + (a.w if dx > 0 else 0)
        y = a.cy + int(dy * abs(x - a.cx) / max(abs(dx), 1))
    else:
        y = a.y + (a.h if dy > 0 else 0)
        x = a.cx + int(dx * abs(y - a.cy) / max(abs(dy), 1))
    return x, y


def center_text(draw: ImageDraw.ImageDraw, box: tuple[int, int, int, int], text: str, font: ImageFont.ImageFont) -> None:
    bbox = draw.textbbox((0, 0), text, font=font)
    tw, th = bbox[2] - bbox[0], bbox[3] - bbox[1]
    x1, y1, x2, y2 = box
    draw.text((x1 + (x2 - x1 - tw) / 2, y1 + (y2 - y1 - th) / 2 - 2), text, font=font, fill="black")


def diamond(n: Node) -> list[tuple[int, int]]:
    return [(n.cx, n.y), (n.x + n.w, n.cy), (n.cx, n.y + n.h), (n.x, n.cy)]


def scaled(n: Node, scale: int) -> Node:
    return Node(n.id, n.value, n.kind, n.x * scale, n.y * scale, n.w * scale, n.h * scale)


def scaled_box(x: int, y: int, w: int, h: int, scale: int) -> tuple[int, int, int, int]:
    return x * scale, y * scale, (x + w) * scale, (y + h) * scale


def render_png() -> None:
    global FONT_TITLE, FONT_ENTITY, FONT_ATTR, FONT_REL, FONT_CARD, FONT_NOTE
    scale = 2
    image = Image.new("RGB", (PAGE_W * scale, PAGE_H * scale), "white")
    draw = ImageDraw.Draw(image)

    FONT_TITLE = load_font(38 * scale, True)
    FONT_ENTITY = load_font(30 * scale, True)
    FONT_ATTR = load_font(20 * scale)
    FONT_REL = load_font(21 * scale, True)
    FONT_CARD = load_font(18 * scale, True)
    FONT_NOTE = load_font(18 * scale)

    nodes = {n.id: scaled(n, scale) for n in [*NODES, *REL_NODES]}

    title_bbox = draw.textbbox((0, 0), TITLE, font=FONT_TITLE)
    title_width = title_bbox[2] - title_bbox[0]
    draw.text(((PAGE_W * scale - title_width) / 2, 35 * scale), TITLE, font=FONT_TITLE, fill="black")

    for attr in ATTRS:
        a, b = nodes[attr.entity], nodes[attr.attr]
        draw.line((edge_point(a, b), edge_point(b, a)), fill="black", width=scale)

    for rel in RELS:
        a, r, b = nodes[rel.a], nodes[rel.rel], nodes[rel.b]
        draw.line((edge_point(a, r), edge_point(r, a)), fill="black", width=2 * scale)
        draw.line((edge_point(r, b), edge_point(b, r)), fill="black", width=2 * scale)

    for n in NODES:
        node = nodes[n.id]
        box = (node.x, node.y, node.x + node.w, node.y + node.h)
        if n.kind == "entity":
            draw.rectangle(box, fill="white", outline="black", width=2 * scale)
            center_text(draw, box, n.value, FONT_ENTITY)
        else:
            draw.ellipse(box, fill="white", outline="black", width=2 * scale)
            center_text(draw, box, n.value, FONT_ATTR)

    for n in REL_NODES:
        node = nodes[n.id]
        pts = diamond(node)
        draw.polygon(pts, fill="white", outline="black")
        draw.line(pts + [pts[0]], fill="black", width=2 * scale)
        center_text(draw, (node.x, node.y, node.x + node.w, node.y + node.h), n.value, FONT_REL)

    for rel in RELS:
        center_text(draw, scaled_box(rel.a_label[0], rel.a_label[1], 28, 24, scale), rel.a_card, FONT_CARD)
        center_text(draw, scaled_box(rel.b_label[0], rel.b_label[1], 28, 24, scale), rel.b_card, FONT_CARD)

    y = 1510
    draw.rectangle(scaled_box(75, y, 58, 26, scale), fill="white", outline="black", width=2 * scale)
    draw.text((145 * scale, y * scale), "实体", font=FONT_NOTE, fill="black")
    draw.ellipse(scaled_box(230, y, 64, 26, scale), fill="white", outline="black", width=2 * scale)
    draw.text((305 * scale, y * scale), "属性", font=FONT_NOTE, fill="black")
    rel_legend = [((405 + 32) * scale, y * scale), ((405 + 64) * scale, (y + 13) * scale), ((405 + 32) * scale, (y + 26) * scale), (405 * scale, (y + 13) * scale)]
    draw.polygon(rel_legend, fill="white", outline="black")
    draw.text((480 * scale, y * scale), "关系；1/N 表示基数", font=FONT_NOTE, fill="black")
    draw.text((75 * scale, 1548 * scale), NOTE, font=FONT_NOTE, fill="black")

    image = image.resize((PAGE_W, PAGE_H), Image.Resampling.LANCZOS)
    image.save(PNG_OUT)
    shutil.copyfile(PNG_OUT, RADIAL_PNG_OUT)


def add_text(root: ET.Element, text_id: str, value: str, x: int, y: int, w: int, h: int, font_size: int = 18, bold: bool = True) -> None:
    cell = ET.SubElement(root, "mxCell", {
        "id": text_id,
        "value": value,
        "style": (
            "text;html=1;strokeColor=none;fillColor=none;align=center;verticalAlign=middle;"
            f"fontSize={font_size};fontStyle={1 if bold else 0};"
        ),
        "vertex": "1",
        "parent": "1",
    })
    ET.SubElement(cell, "mxGeometry", {
        "x": str(x),
        "y": str(y),
        "width": str(w),
        "height": str(h),
        "as": "geometry",
    })


def add_cell(root: ET.Element, n: Node) -> None:
    if n.kind == "entity":
        style = "whiteSpace=wrap;html=1;fontSize=22;fontStyle=1;fillColor=#ffffff;strokeColor=#000000;strokeWidth=2;"
    elif n.kind == "attr":
        style = "ellipse;whiteSpace=wrap;html=1;fontSize=16;fillColor=#ffffff;strokeColor=#000000;strokeWidth=2;"
    else:
        style = "rhombus;whiteSpace=wrap;html=1;fontSize=17;fontStyle=1;fillColor=#ffffff;strokeColor=#000000;strokeWidth=2;"
    cell = ET.SubElement(root, "mxCell", {
        "id": n.id,
        "value": n.value,
        "style": style,
        "vertex": "1",
        "parent": "1",
    })
    ET.SubElement(cell, "mxGeometry", {
        "x": str(n.x),
        "y": str(n.y),
        "width": str(n.w),
        "height": str(n.h),
        "as": "geometry",
    })


def add_edge(root: ET.Element, edge_id: str, source: str, target: str, width: int = 1, label: str = "") -> None:
    cell = ET.SubElement(root, "mxCell", {
        "id": edge_id,
        "value": label,
        "style": (
            f"endArrow=none;html=1;strokeColor=#000000;strokeWidth={width};"
            "fontSize=16;fontStyle=1;labelBackgroundColor=#ffffff;"
        ),
        "edge": "1",
        "parent": "1",
        "source": source,
        "target": target,
    })
    ET.SubElement(cell, "mxGeometry", {"relative": "1", "as": "geometry"})


def write_drawio() -> None:
    mxfile = ET.Element("mxfile", {"host": "app.diagrams.net"})
    diagram = ET.SubElement(mxfile, "diagram", {"id": "er-radial-redrawn", "name": "系统核心业务E-R图"})
    model = ET.SubElement(diagram, "mxGraphModel", {
        "dx": str(PAGE_W),
        "dy": str(PAGE_H),
        "grid": "1",
        "gridSize": "10",
        "guides": "1",
        "tooltips": "1",
        "connect": "1",
        "arrows": "1",
        "fold": "1",
        "page": "1",
        "pageScale": "1",
        "pageWidth": str(PAGE_W),
        "pageHeight": str(PAGE_H),
        "math": "0",
        "shadow": "0",
    })
    root = ET.SubElement(model, "root")
    ET.SubElement(root, "mxCell", {"id": "0"})
    ET.SubElement(root, "mxCell", {"id": "1", "parent": "0"})

    add_text(root, "title", TITLE, 780, 35, 640, 48, font_size=30, bold=True)
    for node in [*NODES, *REL_NODES]:
        add_cell(root, node)
    for i, attr in enumerate(ATTRS, start=1):
        add_edge(root, f"attr_{i}", attr.entity, attr.attr, 1)
    for i, rel in enumerate(RELS, start=1):
        add_edge(root, f"rel_{i}_a", rel.a, rel.rel, 2, rel.a_card)
        add_edge(root, f"rel_{i}_b", rel.rel, rel.b, 2, rel.b_card)

    add_cell(root, Node("legend_entity", "", "entity", 75, 1510, 58, 26))
    add_cell(root, Node("legend_attr", "", "attr", 230, 1510, 64, 26))
    add_cell(root, Node("legend_rel", "", "rel", 405, 1510, 64, 26))
    add_text(root, "legend_entity_text", "实体", 145, 1509, 70, 28, font_size=16, bold=False)
    add_text(root, "legend_attr_text", "属性", 305, 1509, 70, 28, font_size=16, bold=False)
    add_text(root, "legend_rel_text", "关系；1/N 表示基数", 480, 1509, 240, 28, font_size=16, bold=False)
    add_text(root, "note", NOTE, 75, 1545, 1300, 32, font_size=16, bold=False)

    xml_text = ET.tostring(mxfile, encoding="unicode", short_empty_elements=False)
    DRAWIO_OUT.write_text(xml_text, encoding="utf-8")
    RADIAL_DRAWIO_OUT.write_text(xml_text, encoding="utf-8")


def main() -> None:
    ASSET_DIR.mkdir(parents=True, exist_ok=True)
    write_drawio()
    render_png()
    print(DRAWIO_OUT)
    print(PNG_OUT)
    print(RADIAL_DRAWIO_OUT)
    print(RADIAL_PNG_OUT)


if __name__ == "__main__":
    main()
