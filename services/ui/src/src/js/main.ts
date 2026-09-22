// @ts-ignore
import '../css/main.css'

(function () {
    const repoItems = document.querySelectorAll('._repo-item')

    if (!repoItems) {
        return
    }

    repoItems.forEach(item => item.addEventListener('click', handleClick))

    function handleClick(e: Event): void {
        const target = e.target as HTMLElement
        const elem = target.closest('._repo-item')

        if (!elem) {
            return
        }

        const urlAttr = elem.getAttribute('data-url')

        if (!urlAttr) {
            throw new Error('data-url must contain URL')
        }

        location.href = urlAttr
    }
})();
