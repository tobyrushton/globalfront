export const convertTileIdToCoordinates = (tileId: number): { row: number; col: number } => {
    const row = Math.floor(tileId / 200)
    const col = tileId % 200
    return { row, col }
}

export const convertCoordinatesToTileId = (row: number, col: number): number => {
    return row * 200 + col
}