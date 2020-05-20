import { Component, OnInit, OnDestroy } from '@angular/core';
import { BehaviorSubject, Subscription, fromEvent } from 'rxjs';
import { filter, first } from 'rxjs/operators';

@Component({
  selector: 'app-overlay',
  templateUrl: './overlay.component.html',
  styleUrls: ['./overlay.component.css'],
})
export class OverlayComponent implements OnInit, OnDestroy {
  private readonly _visibility: BehaviorSubject<boolean> = new BehaviorSubject(
    false
  );
  readonly visible = this._visibility.asObservable();

  private subscription: Subscription;

  constructor() {}

  ngOnInit() {
    this.subscription = this.visible.subscribe((visible) => {
      if (visible) {
        this.onShow();
      } else {
        this.onHide();
      }
    });
  }

  private onShow(): void {
    fromEvent(document, 'keyup')
      .pipe(
        filter((event) => (event as KeyboardEvent).code === 'Escape'),
        first()
      )
      .subscribe((_) => this.hide());
  }

  private onHide(): void {}

  show(): void {
    this._visibility.next(true);
  }

  hide(): void {
    this._visibility.next(false);
  }

  ngOnDestroy(): void {
    this.subscription.unsubscribe();
  }
}
