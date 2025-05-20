export const drawerItems = [
  {
    title: "Главная",
    icon: "mdi-remote-desktop",
    routeName: "index",
  },
  {
    title: "Все подключения",
    icon: "mdi-monitor-multiple",
    routeName: "connections",
  },
  {
    title: "Новое подключение",
    icon: "mdi-plus-box",
    routeName: "new-connection",
  },
]

export const connectionItems = [
  {
    name: "SSH", // Дефолтное значение для protocol-a
    value: "ssh",
  },
  {
    name: "RDP",
    value: "rdp",
  },
];
