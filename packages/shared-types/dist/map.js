// Helper functions
export function sizeToString(s) {
    return `${s.width}x${s.height}`;
}
export function posToString(p) {
    return `Pos(${p.x},${p.y})`;
}
export function isPosValid(s, p) {
    return p.x >= 1 && p.x <= s.width && p.y >= 1 && p.y <= s.height;
}
export function infoToString(i) {
    return `Info(#${i.id}, ${i.name}, Desc: ${i.desc})`;
}
