import { useEffect, useRef, type RefObject } from 'react'

type Options = {
  root: RefObject<Element | null>
  target: RefObject<Element | null>
  enabled: boolean
  onIntersect: () => void
}

export function useIntersectionObserver({ root, target, enabled, onIntersect }: Options) {
  const callbackRef = useRef(onIntersect)
  callbackRef.current = onIntersect

  useEffect(() => {
    const targetElement = target.current
    if (!enabled || !targetElement) return

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry?.isIntersecting) callbackRef.current()
      },
      { root: root.current, rootMargin: '400px 0px' },
    )
    observer.observe(targetElement)
    return () => observer.disconnect()
  }, [enabled, root, target])
}
