import { ref, onMounted, onBeforeUnmount } from 'vue'

/**
 * Reactive wrapper around window.matchMedia.
 * Returns a ref<boolean> that tracks whether the media query currently matches.
 *
 * The query string is the standard CSS media query syntax, e.g.
 *   useMediaQuery('(max-width: 767px)')   // phone
 *   useMediaQuery('(min-width: 1024px)')  // desktop
 */
export function useMediaQuery(query) {
  const matches = ref(false)
  let mql = null

  const update = () => {
    matches.value = mql.matches
  }

  onMounted(() => {
    mql = window.matchMedia(query)
    update()
    mql.addEventListener('change', update)
  })

  onBeforeUnmount(() => {
    if (mql) {
      mql.removeEventListener('change', update)
    }
  })

  return matches
}