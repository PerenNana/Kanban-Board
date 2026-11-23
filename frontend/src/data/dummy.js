export const dummyBoards = [
  {
    id: 1,
    name: 'Project Alpha',
    columns: [
      {
        id: 1,
        name: 'To Do',
        cards: [
          { id: 1, title: 'Set up project', description: 'Initialize repository and dependencies' },
          { id: 4, title: 'Write README', description: '' },
        ],
      },
      {
        id: 2,
        name: 'In Progress',
        cards: [
          { id: 2, title: 'Develop feature X', description: 'Implement the main feature' },
        ],
      },
      {
        id: 3,
        name: 'Done',
        cards: [
          { id: 3, title: 'Deploy to production', description: 'Release the app' },
        ],
      },
    ],
  },
  {
    id: 2,
    name: 'Uni',
    columns: [
      {
        id: 4,
        name: 'Offen',
        cards: [{ id: 5, title: 'Übung 5', description: 'Vue-Frontend bauen' }],
      },
      {
        id: 5,
        name: 'Erledigt',
        cards: [{ id: 6, title: 'Übung 4', description: 'CI-Pipeline mit Tests' }],
      },
    ],
  },
]