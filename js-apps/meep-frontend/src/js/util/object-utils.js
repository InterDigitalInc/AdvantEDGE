export function deepCopy(source) {
  var dest = JSON.parse(JSON.stringify(source));
  return dest;
}