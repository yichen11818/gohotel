const fs = require("fs");
const path = require("path");
const { execFileSync } = require("child_process");
const pptxgen = require("pptxgenjs");
const JSZip = require("jszip");
const { imageSizingContain } = require("./pptxgenjs_helpers/image");
const {
  warnIfSlideHasOverlaps,
  warnIfSlideElementsOutOfBounds,
} = require("./pptxgenjs_helpers/layout");

const pptx = new pptxgen();
pptx.layout = "LAYOUT_WIDE";
pptx.author = "Codex + PptxGenJS";
pptx.company = "GoHotel Graduation Project";
pptx.subject = "基于 Go 语言与混合推荐算法的智能酒店客房推荐与预约系统设计与实现";
pptx.title = "智能酒店客房推荐与预约系统";
pptx.lang = "zh-CN";
pptx.theme = {
  headFontFace: "Aptos Display",
  bodyFontFace: "Microsoft YaHei UI",
  lang: "zh-CN",
};
pptx.layout = "LAYOUT_WIDE";
pptx.margin = 0;

const ROOT = path.resolve(__dirname, "..", "..");
const OUT_DIR = path.join(ROOT, "output", "ppt", "gohotel-swiss-ikb");
const IMG_DIR = path.join(OUT_DIR, "images");
const OUT_FILE = path.join(OUT_DIR, "gohotel-swiss-ikb.pptx");

const W = 13.333;
const H = 7.5;
const C = {
  paper: "FAFAF8",
  grey1: "F0F0EE",
  grey2: "D4D4D2",
  grey3: "737373",
  ink: "0A0A0A",
  ikb: "002FA7",
  ikbBright: "5B7BFF",
  white: "FFFFFF",
};
const FONT = {
  head: "Aptos Display",
  zhHead: "Microsoft YaHei UI Light",
  body: "Microsoft YaHei UI",
  mono: "Consolas",
};
const ST = pptx.ShapeType;

function img(name) {
  const p = path.join(IMG_DIR, name);
  if (!fs.existsSync(p)) throw new Error(`Missing image: ${p}`);
  return p;
}

function bg(slide, color = C.paper) {
  slide.background = { color };
}

function rect(slide, x, y, w, h, color, opts = {}) {
  slide.addShape(ST.rect, {
    x,
    y,
    w,
    h,
    fill: { color, transparency: opts.transparency ?? 0 },
    line: { color, transparency: opts.lineTransparency ?? 100, width: opts.lineWidth ?? 0 },
  });
}

function line(slide, x, y, w, h = 0, color = C.ink, width = 0.7, opts = {}) {
  slide.addShape(ST.line, {
    x,
    y,
    w,
    h,
    line: { color, width, transparency: opts.transparency ?? 0, dash: opts.dash },
  });
}

function txt(slide, text, x, y, w, h, opts = {}) {
  const props = {
    x,
    y,
    w,
    h,
    fontFace: opts.fontFace || FONT.body,
    fontSize: opts.size ?? 18,
    bold: opts.bold ?? false,
    italic: opts.italic ?? false,
    color: opts.color || C.ink,
    breakLine: opts.breakLine ?? false,
    margin: opts.margin ?? 0,
    valign: opts.valign || "top",
    align: opts.align || "left",
    rotate: opts.rotate,
    charSpace: opts.charSpace,
    paraSpaceAfterPt: opts.paraSpaceAfterPt ?? 0,
    breakLine: opts.breakLine,
    transparency: opts.transparency,
    isTextBox: true,
  };
  if (opts.fit !== false) props.fit = "shrink";
  slide.addText(text, props);
}

function label(slide, text, x, y, w, color = C.grey3, size = 8.5) {
  const hasCjk = /[\u3400-\u9fff]/.test(text);
  txt(slide, hasCjk ? text : text.toUpperCase(), x, y, w, hasCjk ? 0.22 : 0.18, {
    fontFace: hasCjk ? FONT.body : FONT.mono,
    size,
    color,
    charSpace: hasCjk ? 0 : 1.3,
    bold: true,
  });
}

function chrome(slide, left, page, dark = false) {
  const color = dark ? C.white : C.grey3;
  label(slide, left, 0.62, 0.36, 7.2, color, 7.5);
  label(slide, `${String(page).padStart(2, "0")} / 24`, 11.75, 0.36, 0.95, color, 7.5);
}

function foot(slide, text = "GOHOTEL · SWISS IKB DECK", dark = false) {
  const color = dark ? C.white : C.grey3;
  line(slide, 0.62, 6.86, 5.6, 0, color, 0.35, { transparency: 58 });
  label(slide, text, 0.62, 7.04, 4.1, color, 6.8);
  navDots(slide, 6.05, 7.18, dark, Math.max(0, pptx._slides.length - 1));
  txt(slide, "← → 翻页 · B 动态 · ESC 索引", 10.72, 7.05, 2.05, 0.16, {
    fontFace: FONT.body,
    size: 7.2,
    color: dark ? C.ikbBright : C.ikb,
    bold: true,
  });
}

function addNotes(slide, text) {
  slide.addNotes(text);
}

function title(slide, text, sub, page, opts = {}) {
  chrome(slide, opts.chrome || "", page, opts.dark);
  label(slide, opts.kicker || "", 0.72, 1.12, 6.8, opts.dark ? C.ikbBright : C.ikb, 8.5);
  txt(slide, text, 0.72, 1.48, opts.w || 7.5, opts.h || 1.28, {
    fontFace: FONT.head,
    size: opts.size || 34,
    color: opts.dark ? C.white : C.ink,
    bold: false,
    fit: opts.fit,
  });
  if (sub) {
    txt(slide, sub, 0.74, opts.subY || 2.55, opts.subW || 7.2, opts.subH || 0.6, {
      size: opts.subSize || 15.5,
      color: opts.dark ? C.grey2 : C.grey3,
    });
  }
}

function addDotMatrix(slide, x, y, cols, rows, color, size = 0.035, gap = 0.125, tr = 55) {
  for (let r = 0; r < rows; r++) {
    for (let c = 0; c < cols; c++) {
      rect(slide, x + c * gap, y + r * gap, size, size, color, { transparency: tr });
    }
  }
}

function navDots(slide, x, y, dark = false, active = null) {
  const base = dark ? C.white : C.grey3;
  for (let i = 0; i < 24; i++) {
    const color = i === active ? (dark ? C.white : C.ikb) : base;
    rect(slide, x + i * 0.13, y, 0.045, 0.045, color, { transparency: i === active ? 0 : 35 });
  }
}

function addImageContain(slide, file, x, y, w, h, alt) {
  rect(slide, x, y, w, h, C.white);
  slide.addImage({
    path: file,
    ...imageSizingContain(file, x + 0.08, y + 0.08, w - 0.16, h - 0.16),
    altText: alt,
  });
  line(slide, x, y, w, 0, C.ikb, 1.4);
}

function addImageCrop(slide, file, x, y, w, h, alt) {
  rect(slide, x, y, w, h, C.white);
  slide.addImage({
    path: file,
    x,
    y,
    w,
    h,
    altText: alt,
  });
}

function card(slide, x, y, w, h, head, body, opts = {}) {
  rect(slide, x, y, w, h, opts.fill || C.grey1);
  line(slide, x, y, w, 0, opts.accent || C.ikb, opts.lineWidth || 1.6);
  if (opts.no) {
    txt(slide, opts.no, x + 0.16, y + 0.16, 0.55, 0.25, {
      fontFace: FONT.mono,
      size: 10,
      color: opts.noColor || C.ikb,
      bold: true,
    });
  }
  txt(slide, head, x + 0.18, y + (opts.no ? 0.48 : 0.22), w - 0.36, 0.34, {
    size: opts.headSize || 16,
    bold: true,
    color: opts.headColor || C.ink,
  });
  txt(slide, body, x + 0.18, y + (opts.no ? 0.88 : 0.65), w - 0.36, h - (opts.no ? 1.02 : 0.82), {
    size: opts.bodySize || 11.4,
    color: opts.bodyColor || C.grey3,
  });
}

function kpi(slide, x, y, w, num, labelText, opts = {}) {
  line(slide, x, y, w, 0, opts.dark ? C.white : C.ikb, 1.4, { transparency: opts.dark ? 0 : 0 });
  txt(slide, num, x, y + 0.18, w, 0.55, {
    fontFace: FONT.head,
    size: opts.size || 28,
    color: opts.color || (opts.dark ? C.white : C.ink),
    bold: true,
  });
  txt(slide, labelText, x, y + 0.78, w, 0.42, {
    size: opts.labelSize || 10.5,
    color: opts.dark ? C.grey2 : C.grey3,
  });
}

function pill(slide, text, x, y, w, dark = false) {
  rect(slide, x, y, w, 0.28, dark ? C.white : C.ikb);
  txt(slide, text, x + 0.08, y + 0.055, w - 0.16, 0.15, {
    fontFace: FONT.mono,
    size: 6.5,
    color: dark ? C.ikb : C.white,
    bold: true,
    charSpace: 0.7,
  });
}

function imageHero(page, chromeText, imageName, alt, titleText, paragraph, kpis) {
  const slide = pptx.addSlide();
  bg(slide, C.paper);
  addImageContain(slide, img(imageName), 2.78, 0.36, 7.88, 4.08, alt);
  chrome(slide, chromeText, page, false);
  txt(slide, titleText, 0.68, 5.02, 5.2, 0.42, { fontFace: FONT.zhHead, size: 18, color: C.ink });
  txt(slide, paragraph, 0.68, 5.68, 5.9, 0.85, { size: 13.2, color: C.ink });
  kpis.forEach((it, i) => kpi(slide, 6.95 + i * 2.13, 5.08, 1.6, it.num, it.label, { size: 28 }));
  foot(slide);
  addNotes(slide, `${page}. ${titleText}。${paragraph}`);
  return slide;
}

function sectionQuote(page, line1, line2, line3, meta) {
  const slide = pptx.addSlide();
  bg(slide, C.ink);
  addDotMatrix(slide, 9.35, 0.78, 18, 18, C.ikbBright, 0.028, 0.13, 25);
  chrome(slide, "Mid Manifesto · System Design", page, true);
  txt(slide, line1, 0.72, 1.28, 10.2, 0.9, { fontFace: FONT.zhHead, size: 38, color: C.white, fit: false });
  txt(slide, line2, 0.72, 2.26, 10.2, 0.9, { fontFace: FONT.zhHead, size: 38, color: C.white, fit: false });
  txt(slide, line3, 0.72, 3.24, 10.8, 0.9, { fontFace: FONT.zhHead, size: 38, color: C.white, fit: false });
  line(slide, 0.75, 5.22, 4.3, 0, C.ikbBright, 2.2);
  label(slide, meta, 0.76, 5.58, 6, C.grey2, 8);
  foot(slide, "GOHOTEL · DESIGN STATEMENT", true);
  addNotes(slide, "用一句话总结系统设计：推荐可解释、预约状态化、运营数据统一。");
  return slide;
}

function validateSlide(slide) {
  warnIfSlideElementsOutOfBounds(slide, pptx);
  if (process.env.PPTX_FULL_OVERLAP_CHECK === "1") {
    warnIfSlideHasOverlaps(slide, pptx, {
      muteContainment: true,
      ignoreLines: true,
      ignoreDecorativeShapes: true,
    });
  }
}

async function fixPptxNotesSchemaOrder(filePath) {
  const zip = await JSZip.loadAsync(fs.readFileSync(filePath));
  const part = zip.file("ppt/presentation.xml");
  if (!part) return;
  let xml = await part.async("string");
  const notesMatch = xml.match(/<p:notesMasterIdLst>[\s\S]*?<\/p:notesMasterIdLst>/);
  if (notesMatch) {
    // PptxGenJS emits notesMasterIdLst in an order that OfficeCLI's schema
    // validator rejects. PowerPoint/WPS can open it, but removing this list keeps
    // the deck schema-clean while slide notes remain non-critical metadata.
    xml = xml.replace(notesMatch[0], "");
    zip.file("ppt/presentation.xml", xml);
  }
  for (let i = 1; i <= pptx._slides.length; i++) {
    const slidePart = zip.file(`ppt/slides/slide${i}.xml`);
    if (!slidePart) continue;
    let slideXml = await slidePart.async("string");
    if (!slideXml.includes("<p:transition")) {
      slideXml = slideXml.replace(
        /(<p:clrMapOvr>[\s\S]*?<\/p:clrMapOvr>)/,
        '$1<p:transition spd="med"><p:fade/></p:transition>'
      );
      zip.file(`ppt/slides/slide${i}.xml`, slideXml);
    }
  }
  const buf = await zip.generateAsync({
    type: "nodebuffer",
    compression: "DEFLATE",
    compressionOptions: { level: 6 },
  });
  fs.writeFileSync(filePath, buf);
}

function applyNativePowerPointAnimations(filePath) {
  const ps1Path = path.join(OUT_DIR, "apply-gohotel-pptx-animations.ps1");
  const script = String.raw`
param(
  [Parameter(Mandatory = $true)]
  [string]$DeckPath
)

$ErrorActionPreference = "Stop"
$msoAnimEffectAppear = 1
$msoAnimEffectFade = 10
$msoAnimEffectFloat = 30
$msoAnimEffectWipe = 22
$msoAnimTriggerWithPrevious = 2
$msoAnimateLevelNone = 0
$msoAnimDirectionUp = 1
$msoAnimDirectionLeft = 4

function Release-ComObject($Obj) {
  if ($null -ne $Obj) {
    [System.Runtime.InteropServices.Marshal]::ReleaseComObject($Obj) | Out-Null
  }
}

function Test-ShapeHasText($Shape) {
  try {
    return [bool]($Shape.HasTextFrame -and $Shape.TextFrame.HasText)
  } catch {
    return $false
  }
}

function Get-ShapeText($Shape) {
  try {
    if (Test-ShapeHasText $Shape) {
      return [string]$Shape.TextFrame.TextRange.Text
    }
  } catch {}
  return ""
}

function Test-IsStaticDecoration($Shape) {
  $left = [double]$Shape.Left
  $top = [double]$Shape.Top
  $width = [double]$Shape.Width
  $height = [double]$Shape.Height

  if ($width -lt 8 -and $height -lt 8) { return $true }
  if ($top -gt 490) { return $true }
  if ($top -lt 45 -and $height -lt 24) { return $true }
  if ($left -le 5 -and $top -le 5 -and $height -ge 500 -and $width -ge 300) { return $true }
  return $false
}

function Test-IsAnimatableShape($Shape) {
  if (Test-IsStaticDecoration $Shape) { return $false }

  $type = [int]$Shape.Type
  $width = [double]$Shape.Width
  $height = [double]$Shape.Height

  if ($type -eq 17) { return (Test-ShapeHasText $Shape) } # Text box
  if ($type -eq 13) { return $true } # Picture
  if ($type -eq 9) { return ($width -ge 24 -or $height -ge 24) } # Line
  if ($type -eq 1) { return ($width -ge 18 -and $height -ge 8) } # AutoShape
  return $false
}

function Get-EffectKind($Shape) {
  $type = [int]$Shape.Type
  $width = [double]$Shape.Width
  $height = [double]$Shape.Height

  if ($type -eq 17) { return $msoAnimEffectFloat }
  if ($type -eq 9) { return $msoAnimEffectWipe }
  if ($type -eq 1 -and ($width -gt ($height * 2.3) -or $height -gt ($width * 2.3))) { return $msoAnimEffectWipe }
  if ($type -eq 13) { return $msoAnimEffectFade }
  return $msoAnimEffectFade
}

function Add-NativeEntranceEffect($Slide, $Shape, [int]$Order) {
  $sequence = $Slide.TimeLine.MainSequence
  $effectId = Get-EffectKind $Shape

  try {
    $effect = $sequence.AddEffect($Shape, $effectId, $msoAnimateLevelNone, $msoAnimTriggerWithPrevious)
  } catch {
    try {
      $effect = $sequence.AddEffect($Shape, $msoAnimEffectFade, $msoAnimateLevelNone, $msoAnimTriggerWithPrevious)
    } catch {
      $effect = $null
    }
  }

  if ($null -eq $effect) { return $false }

  $delay = [Math]::Min(2.35, 0.06 + ($Order * 0.085))
  $duration = 0.36
  if ([int]$Shape.Type -eq 17) { $duration = 0.44 }
  if ([int]$Shape.Type -eq 13) { $duration = 0.48 }

  try { $effect.Timing.Duration = $duration } catch {}
  try { $effect.Timing.TriggerDelayTime = $delay } catch {}
  try { $effect.Timing.Accelerate = 0.08 } catch {}
  try { $effect.Timing.Decelerate = 0.18 } catch {}

  if ($effectId -eq $msoAnimEffectWipe) {
    try {
      if ([double]$Shape.Width -ge [double]$Shape.Height) {
        $effect.EffectParameters.Direction = $msoAnimDirectionLeft
      } else {
        $effect.EffectParameters.Direction = $msoAnimDirectionUp
      }
    } catch {}
  }

  return $true
}

function Apply-SlideAnimations($Slide) {
  $sequence = $Slide.TimeLine.MainSequence
  for ($i = $sequence.Count; $i -ge 1; $i--) {
    try { $sequence.Item($i).Delete() } catch {}
  }

  $items = New-Object System.Collections.Generic.List[object]
  for ($i = 1; $i -le $Slide.Shapes.Count; $i++) {
    $shape = $Slide.Shapes.Item($i)
    if (Test-IsAnimatableShape $shape) {
      $topBucket = [Math]::Floor(([double]$shape.Top) / 18)
      $items.Add([pscustomobject]@{
        Index = $i
        TopBucket = $topBucket
        Top = [double]$shape.Top
        Left = [double]$shape.Left
        Shape = $shape
        Text = Get-ShapeText $shape
      })
    }
  }

  $ordered = $items | Sort-Object TopBucket, Left, Index
  $order = 0
  foreach ($item in $ordered) {
    if (Add-NativeEntranceEffect $Slide $item.Shape $order) {
      $order++
    }
  }
  return $order
}

$powerPoint = $null
$presentation = $null
try {
  $fullPath = (Resolve-Path -LiteralPath $DeckPath).Path
  $powerPoint = New-Object -ComObject PowerPoint.Application
  $presentation = $powerPoint.Presentations.Open($fullPath, $false, $false, $false)

  $totalEffects = 0
  for ($slideNo = 1; $slideNo -le $presentation.Slides.Count; $slideNo++) {
    $slide = $presentation.Slides.Item($slideNo)
    $totalEffects += Apply-SlideAnimations $slide
  }

  $presentation.Save()
  Write-Host "Applied native PowerPoint object animations: $totalEffects effects across $($presentation.Slides.Count) slides."
} finally {
  if ($null -ne $presentation) {
    try { $presentation.Close() } catch {}
    Release-ComObject $presentation
  }
  if ($null -ne $powerPoint) {
    try { $powerPoint.Quit() } catch {}
    Release-ComObject $powerPoint
  }
  [GC]::Collect()
  [GC]::WaitForPendingFinalizers()
}
`;
  fs.writeFileSync(ps1Path, script, "utf8");
  try {
    execFileSync(
      "powershell.exe",
      ["-NoProfile", "-ExecutionPolicy", "Bypass", "-File", ps1Path, filePath],
      { stdio: "inherit", windowsHide: true }
    );
  } finally {
    if (fs.existsSync(ps1Path)) fs.unlinkSync(ps1Path);
  }
}

// 01 Cover
{
  const slide = pptx.addSlide();
  bg(slide, C.ikb);
  addDotMatrix(slide, 5.35, 0.28, 42, 34, C.white, 0.026, 0.13, 84);
  label(slide, "GOHOTEL · Graduation Project / Swiss IKB", 0.68, 0.43, 4.6, C.white, 7.5);
  label(slide, "SS · 26.05.25 · 01 / 24", 10.8, 0.43, 2.1, C.white, 7.5);
  label(slide, "GO · HYBRID RECOMMENDATION · HOTEL BOOKING", 0.68, 1.03, 5.9, C.white, 7.6);
  txt(slide, "智能酒店客房", 0.66, 2.54, 8.8, 1.05, { fontFace: FONT.zhHead, size: 69, color: C.white, fit: false });
  txt(slide, "推荐与预约系统", 0.66, 3.46, 8.8, 1.05, { fontFace: FONT.zhHead, size: 69, color: C.white, fit: false });
  line(slide, 0.68, 6.15, 12.1, 0, C.white, 0.35, { transparency: 56 });
  txt(slide, "基于 Go 语言、前后端分离架构与混合推荐算法，完成从住客选房、在线预约到后台接待、库存定价、房务工单的完整业务闭环。", 0.68, 6.35, 8.1, 0.48, { size: 13.2, color: C.white });
  label(slide, "3122160125 · 张奕琛 · 赣南科技学院", 0.68, 7.03, 4.7, C.white, 7.2);
  label(slide, "→ SWIPE / ARROW KEYS", 11.05, 6.93, 1.65, C.white, 6.5);
  txt(slide, "← → 翻页 · B 动态 · ESC 索引", 10.72, 7.12, 2.05, 0.16, {
    fontFace: FONT.body,
    size: 7.0,
    color: C.white,
    bold: true,
  });
  navDots(slide, 5.35, 7.27, true, 0);
  addNotes(slide, "开场介绍课题名称、技术路线与项目目标：用 Go、前后端分离和混合推荐算法完成酒店客房推荐与预约闭环。");
}

// 02 Problem
{
  const slide = pptx.addSlide();
  bg(slide, C.ink);
  chrome(slide, "Problem Framing · Business State", 2, true);
  txt(slide, "酒店业务的难点", 0.72, 1.35, 8.5, 0.8, { fontFace: FONT.head, size: 38, color: C.white });
  txt(slide, "不是订房。", 0.72, 2.2, 8.5, 0.82, { fontFace: FONT.head, size: 42, color: C.white });
  rect(slide, 0.68, 3.18, 5.88, 0.88, C.ikb);
  txt(slide, "而是状态同时正确。", 0.84, 3.35, 5.55, 0.46, { fontFace: FONT.head, size: 26, color: C.white });
  ["orders", "rooms", "inventory", "pricing", "housekeeping"].forEach((t, i) => {
    line(slide, 7.15, 1.35 + i * 0.76, 4.55, 0, i === 0 ? C.ikbBright : C.grey3, 1.4);
    txt(slide, t, 7.2, 1.5 + i * 0.76, 3.5, 0.28, { fontFace: FONT.mono, size: 13, color: C.white, bold: i === 0 });
  });
  foot(slide, "BUSINESS STATE CONSISTENCY", true);
  addNotes(slide, "强调酒店系统不是简单订房页面，真正难点是订单、房间、库存、价格和房务状态需要同时保持一致。");
}

// 03 Why now
{
  const slide = pptx.addSlide();
  bg(slide, C.paper);
  title(slide, "Why Now", "中小型酒店需要一套统一的数据底座", 3, { chrome: "Context · Why This System", kicker: "Context" });
  const cols = [
    ["Guest Side", "移动端完成查询、推荐、预约与订单查看", "3 类角色", "住客、前台人员与运营管理员在同一套订单、房间和库存数据上协同。"],
    ["Operation Side", "前台、房务、库存与内容运营交叉发生", "6 条链路", "订单确认、入住退房、房态维护、清洁维修、价格策略和公告活动需要同步更新。"],
    ["Core Goal", "先跑通业务闭环，再加入可解释推荐", "1 个闭环", "从房型浏览到后台处理，系统以可运行、可追踪、可扩展作为毕业设计目标。"],
  ];
  cols.forEach((c, i) => {
    const x = 0.72 + i * 4.12;
    card(slide, x, 3.3, 3.5, 2.38, c[0], c[1], { bodySize: 11.2, headSize: 15.5 });
    txt(slide, c[2], x, 5.96, 2.0, 0.35, { fontFace: FONT.head, size: 20, color: C.ikb, bold: true });
    txt(slide, c[3], x, 6.33, 3.45, 0.43, { size: 9.2, color: C.grey3 });
  });
  foot(slide);
  addNotes(slide, "说明为什么要做这套系统：多角色、多链路、多状态需要统一数据底座。");
}

// 04 Requirements
{
  const slide = pptx.addSlide();
  bg(slide, C.grey1);
  title(slide, "需求不是功能堆叠，而是角色边界", "", 4, { chrome: "Requirement · Functional Boundary", kicker: "Requirement", size: 28, w: 10.5 });
  const items = [
    ["01 / Guest", "住客端", "注册登录、酒店展示、房型筛选、推荐查看、在线预约、订单查询与会员中心。"],
    ["02 / Front Desk", "前台接待", "订单确认、入住办理、退房结算、客人检索与订单状态同步。"],
    ["03 / Room Ops", "房态运维", "房间列表、房型分类、楼层可视化、库存矩阵与房间清洁维修状态。"],
    ["04 / Revenue", "收益与内容", "动态定价规则、活动横幅、公告提醒、文件上传、日志和系统设置。"],
  ];
  items.forEach((it, i) => {
    const x = 0.72 + (i % 2) * 6.06;
    const y = 2.82 + Math.floor(i / 2) * 1.68;
    card(slide, x, y, 5.2, 1.22, it[1], it[2], { no: it[0], fill: C.paper, bodySize: 10.2, headSize: 14.5 });
  });
  foot(slide);
  addNotes(slide, "从四类角色和职责说明项目需求边界，避免把系统讲成无边界的功能堆叠。");
}

// 05 Architecture layers
{
  const slide = pptx.addSlide();
  bg(slide, C.paper);
  chrome(slide, "Architecture · Three Layers", 5, false);
  label(slide, "System Layering", 0.68, 1.02, 3.2, C.grey3, 8.5);
  txt(slide, "三层组织，让多端协同有清晰边界", 0.68, 1.28, 11.2, 0.9, { fontFace: FONT.zhHead, size: 48, color: C.ink, fit: false });
  const layers = [
    ["Presentation", "UniApp + React", "小程序承载住客流程，管理端承载前台与运营流程，二者都通过 HTTP 接口访问后端。"],
    ["Service", "Go + Gin + GORM", "路由接入、JWT 鉴权、Service 业务规则、Repository 数据访问保持分层。"],
    ["Data & Resources", "MySQL / SQLite / COS", "数据库、对象存储、Swagger 文档与时间轮任务提供运行支撑。"],
  ];
  layers.forEach((l, i) => {
    const x = 0.68 + i * 4.18;
    const y = 2.34;
    rect(slide, x, y, 3.95, 4.84, C.grey1);
    label(slide, l[0], x + 0.22, y + 0.23, 2.6, i === 2 ? C.ikb : C.grey3, 8.5);
    txt(slide, l[1], x + 0.22, y + 0.52, 3.25, 0.38, { size: 18, color: C.ink, bold: true });
    rect(slide, x + 2.45, y + 3.32, 0.72, 0.72, C.ikb);
    line(slide, x + 1.36, y + 2.5, 2.58, 0, C.grey2, 0.45);
    line(slide, x + 1.92, y + 2.92, 0, 1.38, C.grey2, 0.45);
    txt(slide, l[2], x + 0.22, y + 4.38, 3.55, 0.36, { size: 8.8, color: C.grey3 });
  });
  foot(slide);
  addNotes(slide, "介绍系统三层：展示层、服务层、数据与资源层。强调前后端分离和后端分层。");
}

// 06 Architecture image
imageHero(6, "Visual Evidence · Overall Architecture", "05-architecture.png", "系统总体架构图", "总体架构证据", "系统采用前后端分离和后端分层架构，小程序、管理端、接口接入、业务服务、数据访问与资源支撑各自承担明确责任。", [
  { num: "2", label: "住客端与管理端" },
  { num: "Go", label: "Gin 接口与 Service 规则" },
  { num: "4", label: "DB、COS、Swagger、时间轮" },
]);

// 07 Module map
{
  const slide = pptx.addSlide();
  bg(slide, C.paper);
  title(slide, "12 Functional Blocks", "完整展示项目模块：从住客入口到运营支撑", 7, { chrome: "Module Map · Functional Coverage", kicker: "Module Map", size: 31, w: 8.8 });
  const blocks = ["账号认证", "酒店展示", "房型浏览", "推荐服务", "在线预约", "订单中心", "前台接待", "房态可视化", "库存矩阵", "动态定价", "工单处理", "上传与日志"];
  blocks.forEach((b, i) => {
    const col = i % 4, row = Math.floor(i / 4);
    const x = 0.72 + col * 3.02;
    const y = 2.55 + row * 1.05;
    rect(slide, x, y, 2.62, 0.74, i === 3 ? C.ikb : C.grey1);
    txt(slide, String(i + 1).padStart(2, "0"), x + 0.14, y + 0.14, 0.42, 0.18, { fontFace: FONT.mono, size: 8.5, color: i === 3 ? C.white : C.ikb, bold: true });
    txt(slide, b, x + 0.68, y + 0.22, 1.65, 0.23, { size: 13, color: i === 3 ? C.white : C.ink, bold: true });
  });
  txt(slide, "Source: thesis chapter 4-5 module design", 0.76, 6.1, 4.0, 0.22, { fontFace: FONT.mono, size: 8, color: C.grey3, charSpace: 0.7 });
  txt(slide, "single hotel / full business loop", 8.68, 6.1, 3.7, 0.22, { fontFace: FONT.mono, size: 8, color: C.ikb, charSpace: 0.7, align: "right" });
  foot(slide);
  addNotes(slide, "列出十二个功能块，展示系统覆盖从住客侧到后台运营侧的完整模块。");
}

// 08 ER image
imageHero(8, "Database · E-R Evidence", "08-er.png", "系统数据库 E-R 图", "数据库 E-R 证据", "数据库围绕用户、房间、订单、库存、价格、工单、公告、行为记录和日志组织，重点保证业务状态可追踪。", [
  { num: "15+", label: "核心实体" },
  { num: "6", label: "订单、房态、清洁等状态" },
  { num: "Service", label: "约束检查集中在服务层" },
]);

// 09 Booking lifecycle
{
  const slide = pptx.addSlide();
  bg(slide, C.paper);
  title(slide, "Linear Flow", "预订与入住：系统中跨模块最多的一条链路", 9, { chrome: "Core Flow · Booking Lifecycle", kicker: "Core Flow" });
  const steps = [
    ["01", "查询房源", "日期与房型条件"],
    ["02", "生成推荐", "行为与特征打分"],
    ["03", "提交预约", "住客信息与日期"],
    ["04", "后台确认", "订单状态更新"],
    ["05", "办理入住", "房间 occupied"],
    ["06", "退房清洁", "dirty 到 clean"],
  ];
  line(slide, 1.08, 3.88, 11.0, 0, C.ikb, 2.2);
  steps.forEach((s, i) => {
    const x = 0.72 + i * 2.02;
    rect(slide, x + 0.25, 3.56, 0.62, 0.62, i === 2 ? C.ikb : C.ink);
    txt(slide, s[0], x + 0.35, 3.75, 0.42, 0.16, { fontFace: FONT.mono, size: 8.5, color: C.white, bold: true, align: "center" });
    txt(slide, s[1], x, 4.46, 1.7, 0.28, { size: 13, bold: true });
    txt(slide, s[2], x, 4.82, 1.74, 0.34, { size: 9.6, color: C.grey3 });
  });
  rect(slide, 1.02, 5.78, 10.6, 0.62, C.grey1);
  txt(slide, "Service 层负责日期校验、库存扣减、房态回写和会员积分更新。", 1.22, 5.96, 10.1, 0.26, { size: 13.5, color: C.ink });
  foot(slide);
  addNotes(slide, "展示预订到入住退房的主流程，并指出 Service 层负责一致性规则。");
}

// 10 Inventory pricing loop
{
  const slide = pptx.addSlide();
  bg(slide, C.grey1);
  title(slide, "Closed Loop", "库存与价格不是静态字段，而是订单流程的一部分", 10, { chrome: "Inventory · Pricing Loop", kicker: "Inventory", size: 32 });
  const loop = [
    ["01 / Query", "查询可售房型", "按房型和日期读取剩余库存。"],
    ["02 / Price", "逐日计算价格", "库存基础价叠加定价规则。"],
    ["03 / Booking", "订单创建扣减", "写入订单后按日期区间更新已订数量。"],
    ["04 / Return", "取消退房回写", "释放库存或恢复房态，进入清洁流程。"],
  ];
  loop.forEach((l, i) => {
    const x = 0.92 + (i % 2) * 5.8;
    const y = 2.55 + Math.floor(i / 2) * 1.62;
    card(slide, x, y, 4.92, 1.18, l[1], l[2], { no: l[0], fill: C.paper, headSize: 14.3, bodySize: 10.4 });
  });
  rect(slide, 2.1, 6.02, 8.9, 0.34, C.ikb);
  txt(slide, "BookingService + InventoryService + PricingRule 共同完成联动。", 2.25, 6.105, 8.6, 0.15, { size: 9.3, color: C.white, fontFace: FONT.mono, bold: true, charSpace: 0.6 });
  foot(slide);
  addNotes(slide, "说明库存与价格要随着查询、下单、取消和退房回写而循环更新。");
}

// 11 Duo compare
{
  const slide = pptx.addSlide();
  bg(slide, C.paper);
  title(slide, "Duo Compare", "住客端追求路径短，管理端追求状态准", 11, { chrome: "Two Ends · Shared Backend", kicker: "Two Ends" });
  line(slide, 6.66, 1.22, 0, 5.38, C.grey2, 0.7);
  card(slide, 0.72, 2.62, 5.35, 2.7, "UniApp Mini Program", "快速找到房，完成预约。首页展示、日期选择、房型列表、详情推荐、确认订单、订单中心和会员服务构成移动端闭环。", { no: "A", fill: C.grey1, headSize: 16, bodySize: 11 });
  card(slide, 7.25, 2.62, 5.35, 2.7, "React Admin", "让后台处理可追踪。运营控制台、订单管理、前台接待、房态可视化、库存矩阵、动态定价和工单处理统一承载。", { no: "B", fill: C.grey1, headSize: 16, bodySize: 11 });
  ["Bearer Token", "行为上报", "订单状态变化"].forEach((t, i) => pill(slide, t, 0.92 + i * 1.55, 5.65, 1.24, false));
  ["订单确认", "库存矩阵", "工单回写"].forEach((t, i) => pill(slide, t, 7.45 + i * 1.55, 5.65, 1.24, false));
  foot(slide);
  addNotes(slide, "对比住客端和管理端，两端使用同一后端，但关注点不同。");
}

// 12 Recommendation signals
{
  const slide = pptx.addSlide();
  bg(slide, C.ink);
  chrome(slide, "Recommendation · Behavior Signals", 12, true);
  label(slide, "Hybrid Signals", 0.72, 1.0, 4, C.ikbBright, 8.5);
  txt(slide, "把用户操作变成可计算的推荐信号", 0.72, 1.38, 8.8, 0.82, { fontFace: FONT.head, size: 32, color: C.white });
  txt(slide, "推荐模块不训练复杂模型，而是把浏览、详情、点击与预订意图转换成权重，再融合协同过滤、房型特征和热度补偿。", 0.76, 2.34, 7.9, 0.8, { size: 13.5, color: C.grey2 });
  const sig = [
    ["view", "浏览房型", "记录基础兴趣入口"],
    ["detail", "查看详情", "提高房型属性匹配权重"],
    ["click", "点击推荐", "反馈推荐结果是否有效"],
    ["book", "预订意图", "作为更强行为信号进入融合打分"],
  ];
  sig.forEach((s, i) => {
    const y = 3.62 + i * 0.62;
    txt(slide, `signal ${String(i + 1).padStart(2, "0")}`, 0.78, y + 0.13, 1.1, 0.12, { fontFace: FONT.mono, size: 6.8, color: C.ikbBright, bold: true });
    line(slide, 2.08, y + 0.18, 1.0 + i * 0.42, 0, C.ikbBright, 3);
    txt(slide, s[0], 4.0, y, 1.0, 0.22, { fontFace: FONT.mono, size: 10, color: C.white, bold: true });
    txt(slide, s[1], 5.15, y, 1.25, 0.22, { size: 12, color: C.white, bold: true });
    txt(slide, s[2], 6.75, y, 4.8, 0.22, { size: 10.3, color: C.grey2 });
  });
  foot(slide, "RECOMMENDATION SIGNAL LEDGER", true);
  addNotes(slide, "解释推荐算法输入：浏览、详情、点击和预订意图。");
}

// 13 Formula
{
  const slide = pptx.addSlide();
  bg(slide, C.paper);
  title(slide, "Recommendation Formula", "轻量混合推荐：有历史行为时融合协同、属性与热度，冷启动时侧重属性和热度。", 13, { chrome: "Algorithm · Score Blend", kicker: "Algorithm", size: 29, w: 10.3, subW: 9.6 });
  rect(slide, 0.72, 3.0, 11.9, 0.95, C.ikb);
  txt(slide, "S = .48C + .37N + .15Pop", 1.0, 3.22, 6.4, 0.4, { fontFace: FONT.mono, size: 21, color: C.white, bold: true });
  txt(slide, "Cold start: S = .68N + .32Pop", 7.65, 3.27, 4.2, 0.25, { fontFace: FONT.mono, size: 11.4, color: C.white, bold: true });
  const weights = [
    ["Collaborative", "0.48", "相似住客带来的协同过滤分"],
    ["Feature", "0.37", "价格、面积、可住人数、折扣率和会员等级"],
    ["Popularity", "0.15", "热门房型与行为热度补偿"],
  ];
  weights.forEach((w, i) => {
    const x = 0.72 + i * 4.05;
    kpi(slide, x, 4.55, 3.25, w[1], w[0], { size: 30, color: i === 0 ? C.ikb : C.ink });
    txt(slide, w[2], x, 5.8, 3.0, 0.38, { size: 10.2, color: C.grey3 });
  });
  txt(slide, "接口返回推荐房型、排序分数和推荐理由，便于解释和调试。", 0.76, 6.48, 8.2, 0.24, { size: 11.5, color: C.grey3 });
  foot(slide);
  addNotes(slide, "给出混合推荐公式和冷启动公式，说明算法轻量、可解释、适合毕业设计演示和调试。");
}

// 14 Guest screenshots
{
  const slide = pptx.addSlide();
  bg(slide, C.paper);
  title(slide, "User Side Evidence", "住客端围绕“看房、订房、查订单”展开", 14, { chrome: "Guest Mini Program · Screens", kicker: "Guest Mini Program", size: 31 });
  const shots = [
    ["13-mini-home-room.png", "01 home / room list", "小程序首页和房型列表"],
    ["14-mini-booking.png", "02 detail / booking", "房间详情和确认订单"],
    ["15-mini-order.png", "03 order status", "订单流程截图"],
  ];
  shots.forEach((s, i) => {
    const x = 0.72 + i * 4.05;
    addImageContain(slide, img(s[0]), x, 2.66, 3.45, 2.48, s[2]);
    label(slide, s[1], x, 5.42, 3.3, C.ikb, 7.5);
  });
  foot(slide);
  addNotes(slide, "展示住客端主要页面：看房、进入详情、提交预约和查看订单状态。");
}

// 15 Booking submit
imageHero(15, "Guest Flow · Booking Submit", "14-mini-booking.png", "小程序房间详情与确认订单界面", "预约提交链路", "小程序端只采集入住人、手机号、身份证、日期和特殊需求；日期合法性、库存可用性、价格计算和库存扣减都在后端完成。", [
  { num: "date", label: "入住与离店日期区间" },
  { num: "stock", label: "房型日期库存" },
  { num: "order", label: "创建订单并扣减库存" },
]);

// 16 Admin dashboard
imageHero(16, "Admin Console · Operation Dashboard", "16-dashboard.png", "管理端运营控制台", "运营控制台", "控制台集中展示入住率、在住房间、待处理工单、订单与房态入口，帮助值班人员快速进入处理任务。", [
  { num: "Order", label: "统一查看和确认订单" },
  { num: "Room", label: "房间、库存与清洁状态" },
  { num: "Todo", label: "工单与运营任务入口" },
]);

// 17 Admin evidence grid
{
  const slide = pptx.addSlide();
  bg(slide, C.paper);
  title(slide, "Management Side Evidence", "管理端承载订单、房态和库存价格三类核心任务", 17, { chrome: "Admin Operations · Evidence Grid", kicker: "Admin Operations", size: 30, w: 10 });
  const shots = [
    ["17-order-manage.png", "01 order management", "订单管理界面"],
    ["18-room-visualization.png", "02 room visualization", "房间可视化界面"],
    ["19-inventory-grid.png", "03 inventory matrix", "库存矩阵界面"],
  ];
  shots.forEach((s, i) => {
    const x = 0.72 + i * 4.05;
    addImageContain(slide, img(s[0]), x, 2.66, 3.45, 2.48, s[2]);
    label(slide, s[1], x, 5.42, 3.3, C.ikb, 7.5);
  });
  foot(slide);
  addNotes(slide, "展示管理端证据图：订单管理、房间可视化和库存矩阵。");
}

// 18 Work orders
imageHero(18, "Housekeeping · Work Orders", "20-repair-orders.png", "维修清洁工单管理界面", "房务工单闭环", "维修工单创建后房间进入 maintenance；维修完成后恢复 available 但标记 dirty；清洁完成后再回到 clean。", [
  { num: "fix", label: "维修工单阻止继续销售" },
  { num: "dirty", label: "退房或维修后等待清洁" },
  { num: "clean", label: "清洁完成后恢复可售条件" },
]);

// 19 Testing
{
  const slide = pptx.addSlide();
  bg(slide, C.ink);
  chrome(slide, "Testing · Coverage Ledger", 19, true);
  label(slide, "Test Scope", 0.72, 1.0, 3.8, C.ikbBright, 8.5);
  txt(slide, "测试围绕后端接口、管理端页面与小程序端流程展开", 0.72, 1.38, 9.1, 0.82, { fontFace: FONT.head, size: 31, color: C.white });
  txt(slide, "测试环境覆盖 Windows 本地开发环境、Go/Gin/GORM 后端、React 管理端与 UniApp H5 预览环境。", 0.76, 2.33, 8.7, 0.5, { size: 13, color: C.grey2 });
  const rows = [
    ["3", "端侧覆盖", "后端接口、管理端页面、小程序端页面"],
    ["9+", "接口范围", "认证、房间、订单、库存、定价、工单、公告、日志等"],
    ["6", "状态流转", "提交、确认、入住、退房、取消、清洁维修"],
    ["1", "业务闭环", "推荐返回、预约提交、库存扣减和后台处理可运行"],
  ];
  rows.forEach((r, i) => {
    const y = 3.28 + i * 0.74;
    line(slide, 0.76, y, 11.6, 0, C.grey3, 0.6, { transparency: 30 });
    txt(slide, r[0], 0.82, y + 0.13, 0.7, 0.33, { fontFace: FONT.head, size: 23, color: C.white, bold: true });
    txt(slide, r[1], 2.0, y + 0.2, 1.45, 0.22, { size: 12, color: C.white, bold: true });
    txt(slide, r[2], 4.0, y + 0.2, 6.8, 0.22, { size: 10.4, color: C.grey2 });
  });
  foot(slide, "TEST COVERAGE LEDGER", true);
  addNotes(slide, "说明测试范围：后端接口、管理端页面、小程序流程和核心业务状态流转。");
}

// 20 Thesis evidence
{
  const slide = pptx.addSlide();
  bg(slide, C.paper);
  title(slide, "Documented System", "论文材料不是空壳：源文档包含章节结构、实现说明、代码片段、表格与截图，支撑本 PPT 的完整项目展示。", 20, { chrome: "Thesis Evidence · Project Scale", kicker: "Thesis Evidence", size: 32, w: 8.8, subW: 9.7 });
  const facts = [
    ["540", "Paragraphs", "正文与代码说明段落"],
    ["28", "Tables", "需求、数据表、状态与测试配置"],
    ["22", "Images", "架构、流程、页面和测试截图"],
    ["29,422", "chars", "指定 `.docx` 论文来源"],
  ];
  facts.forEach((f, i) => {
    const x = 0.72 + i * 3.0;
    kpi(slide, x, 3.58, 2.55, f[0], f[1], { size: i === 3 ? 22 : 32, color: i === 3 ? C.ikb : C.ink });
    txt(slide, f[2], x, 4.92, 2.45, 0.38, { size: 10.1, color: C.grey3 });
  });
  rect(slide, 0.72, 5.88, 11.6, 0.45, C.grey1);
  txt(slide, "内容来源于指定 `.docx` 论文，不凭空扩展为商业 PMS。", 0.94, 6.02, 11.0, 0.18, { size: 10.6, color: C.ink, fontFace: FONT.mono });
  foot(slide);
  addNotes(slide, "用论文统计说明 PPT 内容来自指定文档和真实项目材料。");
}

// 21 Boundary
{
  const slide = pptx.addSlide();
  bg(slide, C.paper);
  title(slide, "What This Project Is / Is Not", "定位清楚，系统才可信", 21, { chrome: "Engineering Position · Boundary", kicker: "Engineering Position", size: 29, w: 10.4 });
  card(slide, 0.72, 2.78, 5.55, 2.72, "No Overclaim", "不是完整商业 PMS。论文明确不把系统表述为复杂商业平台或大规模算法系统，避免超出毕业设计范围。\n\n不覆盖多门店、渠道分销和复杂收益管理；推荐策略不训练深度模型；优惠券入口未接入真实抵扣流程。", { fill: C.grey1, headSize: 16, bodySize: 10.2, accent: C.ink });
  card(slide, 7.0, 2.78, 5.55, 2.72, "Yes Engineering Practice", "是可运行工程实践。重点验证 Go Web、多端协同、混合推荐、库存定价联动与业务状态建模。\n\n核心业务链路可以从用户端跑到后台端；推荐接口返回分数与理由；Service 层集中维护状态一致性。", { fill: C.grey1, headSize: 16, bodySize: 10.2, accent: C.ikb });
  foot(slide);
  addNotes(slide, "明确项目边界：不是商业 PMS，也不是深度学习系统，而是一个完整可运行的工程实践。");
}

// 22 Manifesto
sectionQuote(22, "把推荐做成可解释。", "把预约做成状态机。", "把运营做成同一张数据网。", "GoHotel · design statement");

// 23 Future work
{
  const slide = pptx.addSlide();
  bg(slide, C.grey1);
  title(slide, "Next Iterations", "后续扩展方向：从毕业设计走向真实运营", 23, { chrome: "Future Work · Iteration Notes", kicker: "Future Work", size: 33 });
  const next = [
    ["真实数据训练", "用更长周期行为与订单数据评估推荐效果。"],
    ["优惠券闭环", "将会员资产接入订单抵扣、核销和账务记录。"],
    ["并发一致性", "订单与库存写操作可继续增强事务与锁策略。"],
    ["报表分析", "加入入住率、RevPAR、渠道来源等经营指标。"],
    ["移动端真机", "补充微信小程序真机与多屏尺寸验证。"],
    ["部署运维", "完善云部署、日志告警、备份和权限审计。"],
  ];
  next.forEach((n, i) => {
    const x = 0.72 + (i % 3) * 4.02;
    const y = 2.68 + Math.floor(i / 3) * 1.55;
    card(slide, x, y, 3.42, 1.08, n[0], n[1], { no: String(i + 1).padStart(2, "0"), fill: C.paper, headSize: 13.2, bodySize: 9.6 });
  });
  foot(slide);
  addNotes(slide, "列出后续优化方向，包括真实数据训练、优惠券闭环、并发一致性、报表、真机验证和部署运维。");
}

// 24 Closing
{
  const slide = pptx.addSlide();
  bg(slide, C.paper);
  rect(slide, 0, 0, 6.72, H, C.ikb);
  addDotMatrix(slide, 0.5, 0.42, 32, 28, C.white, 0.025, 0.16, 74);
  label(slide, "24 / 24", 0.48, 0.42, 1.4, C.white, 7.5);
  label(slide, "CLOSING", 5.65, 0.42, 1.0, C.white, 7.5);
  label(slide, "Manifesto", 0.48, 2.24, 1.4, C.white, 7.5);
  txt(slide, "一个系统。", 0.48, 2.7, 4.9, 0.95, { fontFace: FONT.zhHead, size: 56, color: C.white, fit: false });
  txt(slide, "一条闭环。", 0.48, 3.55, 4.9, 0.95, { fontFace: FONT.zhHead, size: 56, color: C.white, fit: false });
  txt(slide, "基于 Go、混合推荐与多端协同，把智能客房推荐和在线预约落成一套可运行的工程实践。", 0.48, 4.85, 4.8, 0.56, { size: 12.8, color: C.white });
  line(slide, 0.48, 6.86, 5.8, 0, C.white, 0.35, { transparency: 56 });
  label(slide, "张奕琛 · Graduation Thesis", 0.48, 7.06, 3.2, C.white, 7.1);
  label(slide, "2026.05.25", 5.4, 7.06, 0.9, C.white, 7.1);
  label(slide, "Takeaways", 7.22, 0.42, 1.35, C.grey3, 7.5);
  label(slide, "03 Rules", 12.28, 0.42, 0.8, C.grey3, 7.5);
  const take = [
    ["01", "业务闭环完整", "覆盖浏览、推荐、预约、确认、入住、退房、库存与房务状态回写。"],
    ["02", "工程边界清楚", "采用前后端分离、后端分层、统一接口与 Service 层一致性约束。"],
    ["03", "推荐能力可解释", "混合策略返回排序结果与推荐理由，适合毕业设计中的演示、调试和答辩说明。"],
  ];
  take.forEach((t, i) => {
    const y = 2.52 + i * 0.98;
    line(slide, 7.22, y, 5.78, 0, i === 2 ? C.ikb : C.grey2, i === 2 ? 1.4 : 0.55);
    txt(slide, t[0], 7.25, y + 0.18, 0.64, 0.34, { fontFace: FONT.head, size: 28, color: i === 2 ? C.ikb : C.ink });
    txt(slide, t[1], 8.2, y + 0.22, 2.2, 0.24, { size: 15.2, color: i === 2 ? C.ikb : C.ink, bold: true });
    txt(slide, t[2], 8.2, y + 0.52, 4.5, 0.26, { size: 9.2, color: C.grey3 });
  });
  label(slide, "→ 完 · GOHOTEL SWISS IKB DECK", 10.45, 7.06, 2.1, C.grey3, 6.5);
  navDots(slide, 5.35, 7.27, false, 23);
  addNotes(slide, "收束三点：业务闭环完整、工程边界清楚、推荐能力可解释。");
}

for (const slide of pptx._slides) {
  validateSlide(slide);
}

fs.mkdirSync(OUT_DIR, { recursive: true });
(async () => {
  const fileName = await pptx.writeFile({ fileName: OUT_FILE, compression: true });
  await fixPptxNotesSchemaOrder(fileName);
  applyNativePowerPointAnimations(fileName);
  await fixPptxNotesSchemaOrder(fileName);
  console.log(`Wrote ${fileName}`);
})();
