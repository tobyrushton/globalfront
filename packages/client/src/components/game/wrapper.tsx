"use client"

import { FC, useState, useRef, PropsWithChildren, MouseEvent, WheelEvent } from "react";

type Coordinates = {
    x: number,
    y: number
}

export const GameWrapper: FC<PropsWithChildren> = ({ children }) => {
    const [position, setPosition] = useState<Coordinates>({ x: 0, y: 0 })
    const [scale, setScale] = useState<number>(1)
    const [dragging, setDragging] = useState<boolean>(false)
    const lastPosition = useRef<Coordinates>({ x: 0, y: 0 })

    const onMouseDown = (e: MouseEvent) => {
        setDragging(true)
        lastPosition.current = { x: e.clientX, y: e.clientY }
    }
    const onMouseUp = (e: MouseEvent) => {
        setDragging(false)
    }

    const handleMouseMove = (e: MouseEvent<HTMLDivElement>) => {
        if (!dragging) return
        const dx = e.clientX - lastPosition.current.x
        const dy = e.clientY - lastPosition.current.y
        setPosition((prev) => ({ x: prev.x + dx, y: prev.y + dy }))
        lastPosition.current = { x: e.clientX, y: e.clientY }
    }

   const handleWheel = (e: WheelEvent<HTMLDivElement>) => {
        const zoomSpeed = 0.005
        const newScale = Math.min(Math.max(0.1, scale + e.deltaY * -zoomSpeed), 5)
        setScale(newScale)
   }

    return (
        <div 
            className="absolute flex w-full h-full justify-center items-center"
            onMouseDown={onMouseDown}
            onMouseUp={onMouseUp}
            onMouseMove={handleMouseMove} 
            onMouseLeave={onMouseUp}
            onWheel={handleWheel}
        >
            <div style={{
                transform: `translate(${position.x}px, ${position.y}px) scale(${scale})`,
                transition: dragging ? 'none' : 'transform 0.1s ease-out'
            }}>
                {children}   
            </div>
        </div>
    )
}