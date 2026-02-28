'use client';

import { useState, useEffect } from 'react';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { useCreateIncome } from '@/hooks/use-budget';
import { SkipForward, AlertTriangle, AlertCircle, Loader2 } from 'lucide-react';
import { format, parseISO } from 'date-fns';
import { useToast } from '@/components/ui/toast';
import { GuestBlockedError } from '@/hooks/use-budget';
import { checkSkippedIncome } from '@/lib/api';
import type { RecurringIncomeWithNextDate } from '@/types/api';

interface SkipOccurrenceDialogProps {
  income: RecurringIncomeWithNextDate | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function SkipOccurrenceDialog({ income, open, onOpenChange }: SkipOccurrenceDialogProps) {
  const createIncome = useCreateIncome();
  const { showToast } = useToast();
  const [skipDate, setSkipDate] = useState('');
  const [reason, setReason] = useState('');
  const [isChecking, setIsChecking] = useState(false);
  const [alreadySkipped, setAlreadySkipped] = useState(false);

  useEffect(() => {
    if (open && skipDate) {
      checkForExistingSkip();
    }
  }, [skipDate, open]);

  const checkForExistingSkip = async () => {
    if (!skipDate) return;
    setIsChecking(true);
    try {
      const isSkipped = await checkSkippedIncome(skipDate);
      setAlreadySkipped(isSkipped);
    } catch (error) {
      console.error('Error checking skip:', error);
    } finally {
      setIsChecking(false);
    }
  };

  const handleOpen = (isOpen: boolean) => {
    if (isOpen && income) {
      setSkipDate(income.next_occurrence);
      setReason('');
      setAlreadySkipped(false);
    }
    onOpenChange(isOpen);
  };

  const handleSkip = async () => {
    if (!income || alreadySkipped) return;

    try {
      await createIncome.mutateAsync({
        amount: -income.amount,
        currency: income.currency,
        date: skipDate,
        description: income.description 
          ? `Skipped: ${income.description}` 
          : 'Skipped income',
        notes: reason ? `Skip reason: ${reason}` : undefined,
      });
      showToast(`Skipped ${income.currency} ${income.amount.toFixed(2)} for ${format(parseISO(skipDate), 'MMM d, yyyy')}`, 'success');
      handleOpen(false);
    } catch (error) {
      if (error instanceof GuestBlockedError) {
        showToast(error.message, 'warning');
      } else {
        showToast('Failed to skip occurrence', 'error');
      }
    }
  };

  if (!income) return null;

  return (
    <Dialog open={open} onOpenChange={handleOpen}>
      <DialogContent className="sm:max-w-[450px]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <SkipForward className="h-5 w-5" />
            Skip Recurring Income
          </DialogTitle>
          <DialogDescription>
            Skip one occurrence of your recurring income. This will create a negative income entry.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 py-4">
          <div className="p-4 rounded-lg bg-muted space-y-2">
            <div className="flex justify-between">
              <span className="text-muted-foreground">Description</span>
              <span className="font-medium">{income.description || 'Income'}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Amount</span>
              <span className="font-medium text-red-600">
                -{income.currency} {income.amount.toFixed(2)}
              </span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Frequency</span>
              <span className="font-medium capitalize">{income.recurring_type || 'N/A'}</span>
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="skip-date">Skip Date</Label>
            <Input
              id="skip-date"
              type="date"
              value={skipDate}
              onChange={(e) => {
                setSkipDate(e.target.value);
                setAlreadySkipped(false);
              }}
            />
            {isChecking && (
              <div className="flex items-center gap-2 text-xs text-muted-foreground">
                <Loader2 className="h-3 w-3 animate-spin" />
                Checking for existing skip...
              </div>
            )}
            {alreadySkipped && !isChecking && (
              <div className="flex items-center gap-2 text-xs text-red-500">
                <AlertCircle className="h-3 w-3" />
                This date has already been skipped
              </div>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="reason">Reason (optional)</Label>
            <Textarea
              id="reason"
              placeholder="e.g., vacation, sick leave, holiday..."
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              rows={2}
            />
          </div>

          <div className="flex items-start gap-2 p-3 rounded-lg bg-amber-500/10 text-amber-700 dark:text-amber-400">
            <AlertTriangle className="h-4 w-4 mt-0.5 shrink-0" />
            <p className="text-xs">
              This will create a negative income entry. You can delete it later to undo the skip.
            </p>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => handleOpen(false)}>
            Cancel
          </Button>
          <Button 
            variant="destructive" 
            onClick={handleSkip}
            disabled={createIncome.isPending || !skipDate || alreadySkipped || isChecking}
          >
            {createIncome.isPending ? 'Skipping...' : 'Skip Income'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
