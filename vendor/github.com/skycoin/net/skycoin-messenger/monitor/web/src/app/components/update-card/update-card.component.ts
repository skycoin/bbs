import { Component, OnInit, ViewEncapsulation, Inject } from '@angular/core';
import { ApiService } from '../../service/api/api.service';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';

const NOUPGRADE = 'No Upgrade Available';
const UPGRADE = 'Upgrade Available';

@Component({
  selector: 'app-update-card',
  templateUrl: './update-card.component.html',
  styleUrls: ['./update-card.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class UpdateCardComponent implements OnInit {
  progressValue = 0;
  progressTask = null;
  updateStatus = NOUPGRADE;
  hasUpdate = false;
  nodeUrl = '';
  dialogRef: MatDialogRef<UpdateCardComponent>;
  constructor(private api: ApiService, @Inject(MAT_DIALOG_DATA) public data: { version?: string, tag?: string }) { }

  ngOnInit() {
    this.api.checkUpdate(this.data.tag, this.data.version).subscribe((res: Update) => {
      this.hasUpdate = res.Update;
      if (this.hasUpdate) {
        this.updateStatus = UPGRADE;
      } else {
        this.updateStatus = NOUPGRADE;
      }
    });
  }
  startDownload(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.progressTask = setInterval(() => {
      this.progressValue += 1;
      if (this.progressValue >= 100) {
        clearInterval(this.progressTask);
      }
    }, 100);
    this.api.runNodeupdate(this.nodeUrl).subscribe(result => {
      this.dialogRef.close(result);
    });
  }
  getUpgradeStatus() {
    return false;
  }
}

export interface Update {
  Force?: boolean;
  Update?: boolean;
  Latest?: string;
}
