import { animate, AnimationEntryMetadata, state, style, transition, trigger } from '@angular/core';
export const DialogAnimation: AnimationEntryMetadata =
  trigger('dialogInOut', [
    state('void', style({ opacity: 1, transform: 'scale3d(1, 1, 1)' })),
    transition('void => *', [
      style({
        opacity: 0,
        transform: 'scale3d(.1, .1, .1)'
      }),
      animate(250)
    ]),
    transition('* => void', [
      animate(250, style({
        opacity: 0,
        transform: 'scale3d(.1, .1, .1)'
      }))
    ])
  ])
