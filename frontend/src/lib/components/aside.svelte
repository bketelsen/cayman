<script lang="ts">
    import {
        CpuIcon,
        TreePalm,
        House,
        Logs,
        Cylinder,
        Container,
        TrendingUpDown,
        ServerIcon,
    } from "@lucide/svelte";
    import { page } from "$app/state";

</script>

<aside class="drawer-side z-10">
    <label for="my-drawer" class="drawer-overlay"></label>
    <!-- sidebar menu -->
    <nav
        class="bg-base-100 flex min-h-screen w-72 flex-col gap-2 overflow-y-auto px-6 py-10"
    >
        <div class="mx-4 flex items-center gap-2 font-black">
            <TreePalm class="h-8 w-8" />
            Cayman
        </div>
        <ul class="menu text-base-content bg-base-200 w-full">
            <li>
                <a
                    href="/"
                    class={page.url.pathname === "/" ? "menu-active" : ""}
                    aria-current={page.url.pathname === "/"}
                >
                    <House class="h-5 w-5" />
                    Dashboard
                </a>
            </li>
            <li>
                <a
                    href="/system"
                    class={page.url.pathname === "/system" ? "menu-active" : ""}
                    aria-current={page.url.pathname === "/system"}
                >
                    <ServerIcon class="h-5 w-5" />
                    System
                </a>
            </li>
            <li>
                <a
                    href="/host"
                    class={page.url.pathname === "/host" ? "menu-active" : ""}
                    aria-current={page.url.pathname === "/host"}
                >
                    <CpuIcon class="h-5 w-5" />
                    Host
                </a>
            </li>
            <li>
                <a
                    href="/storage"
                    aria-disabled="true"
                    class={page.url.pathname === "/storage"
                        ? "menu-active"
                        : ""}
                    aria-current={page.url.pathname === "/storage"}
                >
                    <Cylinder class="h-5 w-5" />
                    Storage
                </a>
            </li>

            <li>
                <details>
                    <summary>
                        <Container class="h-5 w-5" />
                        Virtualization
                    </summary>
                    <ul>
                        <li>
                            <a
                                href="/virt/kvm"
                                class={page.url.pathname === "/virt/kvm"
                                    ? "menu-active"
                                    : ""}
                                aria-current={page.url.pathname === "/virt/kvm"}
                                >KVM</a
                            >
                        </li>

                        <li>
                            <a
                                href="/virt/podman"
                                class={page.url.pathname === "/virt/podman"
                                    ? "menu-active"
                                    : ""}
                                aria-current={page.url.pathname ===
                                    "/virt/podman"}>Podman</a
                            >
                        </li>
                        <li>
                            <a
                                href="/virt/docker"
                                class={page.url.pathname === "/virt/docker"
                                    ? "menu-active"
                                    : ""}
                                aria-current={page.url.pathname ===
                                    "/virt/docker"}>Docker</a
                            >
                        </li>
                        <li>
                            <a
                                href="/virt/incus"
                                class={page.url.pathname === "/virt/incus"
                                    ? "menu-active"
                                    : ""}
                                aria-current={page.url.pathname ===
                                    "/virt/incus"}>Incus</a
                            >
                        </li>
                    </ul>
                </details>
            </li>
            <li>
                <a
                    href="/logs"
                    class={page.url.pathname === "/logs"
                        ? "menu-active"
                        : ""}
                        aria-current={page.url.pathname === "/logs"}
                    >
                        <Logs class="h-5 w-5" />
                        Logs
                    </a>

            </li>
            <li>
                <a
                    href="/metrics"
                    class={page.url.pathname === "/metrics"
                        ? "menu-active"
                        : ""}
                    aria-current={page.url.pathname === "/metrics"}
                >
                    <TrendingUpDown class="h-5 w-5" />
                    Metrics
                </a>
            </li>
        </ul>
    </nav>
    <!-- /sidebar menu -->
</aside>
