export function paginate(items, currentPage, pageSize) {
    if (!items || items.length === 0) return [];
    const startIndex = (currentPage - 1) * pageSize;
    return items.slice(startIndex, startIndex + pageSize);
}