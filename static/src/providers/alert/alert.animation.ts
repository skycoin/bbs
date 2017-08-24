import { animate, AnimationEntryMetadata, state, style, transition, trigger } from '@angular/core';

export const AlertAnimation: AnimationEntryMetadata =
  trigger('alertInOut', [
    state('void', style({ opacity: 1, transform: 'translateX(0)' })),
    transition('void => *', [
      style({
        opacity: 0,
        transform: 'translateX(100%)'
      }),
      animate('0.3s ease-in')
    ]),
    transition('* => void', [
      animate('0.3s 0.1s ease-out', style({
        opacity: 0,
        transform: 'translateX(100%) scale(.5)'
      }))
    ])
  ])
export const DialogAnimation: AnimationEntryMetadata =
  trigger('dialogInOut', [
    state('void', style({ opacity: 1, transform: 'scale3d(1)' })),
    transition('void => *', [
      style({
        opacity: 0,
        transform: 'scale3d(0.5)'
      }),
      animate('0.3s ease-in')
    ]),
    transition('* => void', [
      animate('0.3s 0.1s ease-out', style({
        opacity: 0,
        transform: 'scale3d(0.5)'
      }))
    ])
  ])
